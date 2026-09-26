const Redis = require('ioredis');
const { spawn } = require('child_process');
const fs = require('fs').promises;
const path = require('path');

const redis = new Redis();

const REPO_ROOT = process.env.REPO_ROOT || path.join(process.env.HOME, 'Documents/LUKMAN/software-engineering-lab');
const AGENT_ROOT = process.env.AGENT_ROOT || path.join(process.env.HOME, 'research-agent');
const PROMPTS = path.join(AGENT_ROOT, 'prompts');
const TASK_ROOT = path.join(AGENT_ROOT, 'tasks');
const TASK_PENDING = path.join(TASK_ROOT, 'pending');
const TASK_DONE = path.join(TASK_ROOT, 'done');
const LOG_ROOT = path.join(AGENT_ROOT, 'logs');

const MAX_RESEARCH_REVISIONS = 6;
const MAX_ENGINEERING_REVISIONS = 6;
const MAX_WRITER_REVISIONS = 6;

const MODEL_DEFAULT = "9router/free-combo";
const MODEL_CRITICAL = "9router/ag-combo";
const MODEL_FALLBACK = "9router/ag-combo";
const MODEL_AUDITOR_OS = "9router/auditor-opensource-combo";

async function runOpencode(lab, stage, promptFile, instruction, model = MODEL_DEFAULT) {
    const logDir = path.join(LOG_ROOT, path.basename(lab));
    await fs.mkdir(logDir, { recursive: true });
    const logFile = path.join(logDir, `${stage}.log`);

    const promptContent = await fs.readFile(promptFile, 'utf8');
    const input = `${promptContent}\n---\n${instruction}`;

    async function execute(mod) {
        return new Promise((resolve, reject) => {
            const args = ['run'];
            if (mod) args.push('-m', mod);
            
            console.log(`[${new Date().toISOString()}] START: ${stage} | LAB: ${lab} | Model: ${mod}`);
            const child = spawn('opencode', args, { cwd: REPO_ROOT });
            
            child.stdout.on('data', data => fs.appendFile(logFile, data).catch(()=>{}));
            child.stderr.on('data', data => fs.appendFile(logFile, data).catch(()=>{}));

            child.on('close', code => resolve(code));
            child.on('error', err => reject(err));

            child.stdin.write(input);
            child.stdin.end();
        });
    }

    let code = await execute(model);
    if (code !== 0 && model !== MODEL_FALLBACK) {
        console.warn(`WARNING: Model ${model} failed (exit ${code}). Fallback to ${MODEL_FALLBACK}...`);
        code = await execute(MODEL_FALLBACK);
    }

    if (code !== 0) {
        throw new Error(`Stage ${stage} failed with exit code ${code}`);
    }
    
    console.log(`[${new Date().toISOString()}] DONE: ${stage}`);
}

async function parseVerdict(file) {
    try {
        const content = await fs.readFile(file, 'utf8');
        const match = content.match(/(APPROVED_WITH_WARNINGS|NEEDS_REVISION|REJECTED|APPROVED)/g);
        if (match && match.length > 0) {
            return match[match.length - 1];
        }
    } catch (err) {
        // Ignore read error
    }
    return 'NEEDS_REVISION';
}

async function researchPipeline(lab, taskFile) {
    let revisionCount = 0;
    const taskPath = path.join(TASK_PENDING, taskFile);
    const taskContent = await fs.readFile(taskPath, 'utf8');

    while (true) {
        if (revisionCount === 0) {
            await runOpencode(lab, "01-research", path.join(PROMPTS, "researcher.md"), 
`Target lab: ${lab}\n\nExecute the topic specification below for the RESEARCH PHASE ONLY.\n\nPIPELINE BOUNDARY:\n- Research only.\n- Do not implement source code.\n- Do not write tests.\n- Do not generate publication content.\n- Engineering is handled by a separate Engineer Agent.\n- Write research output only inside: ${lab}/research/\n\nTOPIC SPECIFICATION:\n\n${taskContent}`, MODEL_DEFAULT);
        } else {
            await runOpencode(lab, `03-research-revision-r${revisionCount}`, path.join(PROMPTS, "research-reviser.md"),
`Address the audit findings for target lab:\n\n${lab}\n\nPIPELINE OVERRIDE:\n- Revise research only.\n- Do not write code.\n- Write revision records to:\n  ${lab}/research-revision/\n- Finish with READY_FOR_RESEARCH_REAUDIT.`, MODEL_DEFAULT);
        }

        const auditRound = revisionCount + 1;
        await fs.rm(path.join(REPO_ROOT, lab, 'research-audit'), { recursive: true, force: true });
        
        await runOpencode(lab, `02-research-audit-r${auditRound}`, path.join(PROMPTS, "research-auditor.md"),
`Audit the research for target lab:\n\n${lab}\n\nPIPELINE OVERRIDE:\n- Audit research only.\n- Do not audit implementation/code in this stage.\n- Do not modify research files.\n- Write all audit output to: ${lab}/research-audit/\n- Final verdict must be written to:\n  ${lab}/research-audit/07-verdict.md`, MODEL_CRITICAL);

        const verdict = await parseVerdict(path.join(REPO_ROOT, lab, 'research-audit/07-verdict.md'));
        console.log(`Research verdict: ${verdict}`);

        if (verdict === 'APPROVED') {
            return true;
        } else if (verdict === 'NEEDS_REVISION' || verdict === 'APPROVED_WITH_WARNINGS') {
            if (revisionCount >= MAX_RESEARCH_REVISIONS) {
                console.log("BLOCKED: maximum research revisions reached.");
                throw new Error("Max revisions reached in research.");
            }
            revisionCount++;
        } else {
            console.log(`BLOCKED: Research verdict was ${verdict}.`);
            throw new Error(`Research blocked by verdict: ${verdict}`);
        }
    }
}

async function engineeringPipeline(lab) {
    let revisionCount = 0;
    while (true) {
        if (revisionCount === 0) {
            await runOpencode(lab, "04-engineering", path.join(PROMPTS, "engineer.md"),
`Implement the approved technical lab:\n\n${lab}\n\nPIPELINE BOUNDARY:\n- Write and edit source code.\n- Write tests.\n- Do not generate publication content.`, MODEL_CRITICAL);
        } else {
            await runOpencode(lab, `05-engineering-revision-r${revisionCount}`, path.join(PROMPTS, "engineer-reviser.md"),
`Address the audit findings for target lab:\n\n${lab}\n\nPIPELINE OVERRIDE:\n- Revise implementation and tests only.\n- Do not generate publication content.\n- Write revision records to:\n  ${lab}/engineering-revision/\n- Finish with READY_FOR_ENGINEERING_REAUDIT.`, MODEL_CRITICAL);
        }

        const auditRound = revisionCount + 1;
        await fs.rm(path.join(REPO_ROOT, lab, 'engineering-audit'), { recursive: true, force: true });
        await fs.rm(path.join(REPO_ROOT, lab, 'engineering-audit-opensource'), { recursive: true, force: true });

        // Run both auditors asynchronously in parallel
        await Promise.all([
            runOpencode(lab, `06-engineering-audit-r${auditRound}`, path.join(PROMPTS, "engineering-auditor.md"),
`Audit the implementation for target lab:\n\n${lab}\n\nPIPELINE OVERRIDE:\n- Audit implementation and tests only.\n- Do not audit research/content in this stage.\n- Do not modify code.\n- Write all audit output to: ${lab}/engineering-audit/\n- Final verdict must be written to:\n  ${lab}/engineering-audit/06-verdict.md`, MODEL_CRITICAL),
            
            runOpencode(lab, `06-engineering-audit-opensource-r${auditRound}`, path.join(PROMPTS, "engineering-auditor.md"),
`Audit the implementation for target lab:\n\n${lab}\n\nPIPELINE OVERRIDE:\n- Audit implementation and tests only.\n- Do not audit research/content in this stage.\n- Do not modify code.\n- Write all audit output to: ${lab}/engineering-audit-opensource/\n- Final verdict must be written to:\n  ${lab}/engineering-audit-opensource/06-verdict.md`, MODEL_AUDITOR_OS)
        ]);

        const verdict1 = await parseVerdict(path.join(REPO_ROOT, lab, 'engineering-audit/06-verdict.md'));
        const verdict2 = await parseVerdict(path.join(REPO_ROOT, lab, 'engineering-audit-opensource/06-verdict.md'));
        
        console.log(`Engineering verdict 1 (default): ${verdict1}`);
        console.log(`Engineering verdict 2 (opensource): ${verdict2}`);

        let finalVerdict = "NEEDS_REVISION";
        if (verdict1 === 'APPROVED' && verdict2 === 'APPROVED') {
            finalVerdict = 'APPROVED';
        } else if (verdict1 === 'REJECTED' || verdict2 === 'REJECTED') {
            finalVerdict = 'REJECTED';
        } else if (verdict1 === 'APPROVED_WITH_WARNINGS' || verdict2 === 'APPROVED_WITH_WARNINGS') {
            finalVerdict = 'APPROVED_WITH_WARNINGS';
        }

        if (finalVerdict === 'APPROVED') {
            return true;
        } else if (finalVerdict === 'NEEDS_REVISION' || finalVerdict === 'APPROVED_WITH_WARNINGS') {
            if (revisionCount >= MAX_ENGINEERING_REVISIONS) {
                console.log("BLOCKED: maximum engineering revisions reached.");
                throw new Error("Max revisions reached in engineering.");
            }
            revisionCount++;
        } else {
            console.log(`BLOCKED: Engineering verdict was ${finalVerdict}.`);
            throw new Error(`Engineering blocked by verdict: ${finalVerdict}`);
        }
    }
}

async function writerPipeline(lab) {
    let revisionCount = 0;
    while (true) {
        if (revisionCount === 0) {
            await runOpencode(lab, "07-technical-writer", path.join(PROMPTS, "technical-writer.md"),
`Draft the final publication content for target lab:\n\n${lab}\n\nPIPELINE BOUNDARY:\n- Draft technical publication only.\n- Do not edit research or code.\n- Output strictly inside: ${lab}/content/\n- Finish with READY_FOR_CONTENT_ADAPTER.`, MODEL_DEFAULT);
        } else {
            await runOpencode(lab, `09-content-revision-r${revisionCount}`, path.join(PROMPTS, "technical-writer-reviser.md"),
`Address the audit findings for target lab content:\n\n${lab}\n\nPIPELINE OVERRIDE:\n- Revise content only.\n- Do not edit research or code.\n- Output strictly inside: ${lab}/content/\n- Finish with READY_FOR_CONTENT_ADAPTER.`, MODEL_DEFAULT);
        }

        const auditRound = revisionCount + 1;
        await fs.rm(path.join(REPO_ROOT, lab, 'content-audit'), { recursive: true, force: true });
        
        await runOpencode(lab, `08-content-audit-r${auditRound}`, path.join(PROMPTS, "technical-writer-auditor.md"),
`Audit the technical publication content for target lab:\n\n${lab}\n\nPIPELINE OVERRIDE:\n- Audit content only.\n- Do not audit research/code.\n- Do not modify files.\n- Write all audit output to: ${lab}/content-audit/\n- Final verdict must be written to:\n  ${lab}/content-audit/09-verdict.md`, MODEL_DEFAULT);

        const verdict = await parseVerdict(path.join(REPO_ROOT, lab, 'content-audit/09-verdict.md'));
        console.log(`Writer verdict: ${verdict}`);

        if (verdict === 'APPROVED') {
            return true;
        } else if (verdict === 'NEEDS_REVISION' || verdict === 'APPROVED_WITH_WARNINGS') {
            if (revisionCount >= MAX_WRITER_REVISIONS) {
                console.log("BLOCKED: maximum writer revisions reached.");
                throw new Error("Max revisions reached in writer.");
            }
            revisionCount++;
        } else {
            console.log(`BLOCKED: Writer verdict was ${verdict}.`);
            throw new Error(`Writer blocked by verdict: ${verdict}`);
        }
    }
}

async function resolveLabFromTask(taskPath) {
    try {
        const content = await fs.readFile(taskPath, 'utf8');
        const match = content.match(/labs\/[a-zA-Z0-9_-]+/);
        return match ? match[0] : null;
    } catch {
        return null;
    }
}

async function archiveTaskDone(taskFile) {
    const source = path.join(TASK_PENDING, taskFile);
    let target = path.join(TASK_DONE, taskFile);
    
    try {
        await fs.access(target);
        const stamp = new Date().toISOString().replace(/[:.]/g, '-');
        target = path.join(TASK_DONE, `${taskFile.replace('.md', '')}-${stamp}.md`);
    } catch (e) {
        // Target doesn't exist, which is good
    }
    
    await fs.rename(source, target);
    console.log(`TASK DONE: ${target}`);
}

async function processTask(taskFile) {
    const lockKey = `lock:task:${taskFile}`;
    const acquired = await redis.set(lockKey, "1", "NX", "EX", 3600); // 1 hour lock
    
    if (!acquired) {
        console.log(`SKIP: ${taskFile} is being processed (locked in Redis).`);
        return;
    }

    try {
        const isDone = await redis.sismember("lab:tasks:done", taskFile);
        if (isDone) {
            console.log(`SKIP: ${taskFile} is already completed (in lab:tasks:done).`);
            return;
        }

        await redis.sadd("lab:tasks:processing", taskFile);

        const taskPath = path.join(TASK_PENDING, taskFile);
        const lab = await resolveLabFromTask(taskPath);

        if (!lab) {
            throw new Error(`Could not resolve lab path from task: ${taskFile}`);
        }

        console.log(`Processing ${taskFile} -> ${lab}`);
        
        await researchPipeline(lab, taskFile);
        await engineeringPipeline(lab);
        await writerPipeline(lab);
        await archiveTaskDone(taskFile);

        await redis.sadd("lab:tasks:done", taskFile);
        
        console.log(`\nCOMPLETED: ${taskFile} -> ${lab}`);

    } catch (err) {
        console.error(`ERROR processing ${taskFile}:`, err.message);
    } finally {
        await redis.srem("lab:tasks:processing", taskFile);
        await redis.del(lockKey);
    }
}

async function main() {
    await fs.mkdir(TASK_PENDING, { recursive: true });
    await fs.mkdir(TASK_DONE, { recursive: true });
    await fs.mkdir(LOG_ROOT, { recursive: true });

    const files = await fs.readdir(TASK_PENDING);
    const pendingTasks = files.filter(f => f.startsWith('lab-') && f.endsWith('-topic.md'));

    if (pendingTasks.length === 0) {
        console.log("No pending tasks.");
        process.exit(0);
    }

    console.log(`Found ${pendingTasks.length} tasks.`);
    
    // Process all tasks concurrently!
    await Promise.all(pendingTasks.map(task => processTask(task)));
    
    console.log("All concurrent task loops finished.");
    process.exit(0);
}

main().catch(err => {
    console.error(err);
    process.exit(1);
});
