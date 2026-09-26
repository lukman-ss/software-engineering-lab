# Evidence

## Evidence 1

Claim: Blue-Green deployment enables rapid rollback and cut-overs with zero or minimal downtime by switching routers between two identical environments.

Evidence: "The blue-green deployment approach does this by ensuring you have two production environments, as identical as possible. At any time one of them, let's say blue for the example, is live. As you prepare a new release of your software you do your final stage of testing in the green environment. Once the software is working in the green environment, you switch the router so that all incoming requests go to the green environment - the blue one is now idle. Blue-green deployment also gives you a rapid way to rollback - if anything goes wrong you switch the router back to your blue environment."

Source: Blue Green Deployment
URL: https://martinfowler.com/bliki/BlueGreenDeployment.html
Confidence: HIGH
Corroborated By: Jez Humble and Dave Farley (Continuous Delivery book)
Notes: Martin Fowler notes that schema changes complicate this, requiring database changes to be decoupled from application upgrades.

---

## Evidence 2

Claim: Database schema changes during zero-downtime deployment require breaking the change into three distinct phases: Expand, Migrate, and Contract (Parallel Change).

Evidence: "Parallel change, also known as expand and contract, is a pattern to implement backward-incompatible changes to an interface in a safe manner, by breaking the change into three distinct phases: expand, migrate, and contract... Most database refactorings follow the parallel change pattern, where the migrate phase is the transition period between the original and the new schema, until all database access code has been updated to work with the new schema."

Source: Parallel Change
URL: https://martinfowler.com/bliki/ParallelChange.html
Confidence: HIGH
Corroborated By: Martin Fowler (Blue Green Deployment, 2010 update)
Notes: Essential for zero-downtime deployment when two application versions run concurrently during rollout.

---

## Evidence 3

Claim: Readiness probes determine when a container is ready to accept traffic; failing readiness detaches the container from load balancer endpoints without restarting it. Liveness probes determine when to restart an unhealthy container.

Evidence: "Liveness probes determine when to restart a container. For example, liveness probes could catch a deadlock, where an application is running, but unable to make progress... Readiness probes determine when a container is ready to accept traffic... If the readiness probe returns a failed state, the EndpointSlice controller removes the Pod's IP address from the EndpointSlices of all Services that match the Pod."

Source: Pod Lifecycle
URL: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/
Confidence: HIGH
Corroborated By: Kubernetes official best practices
Notes: Marking an instance ready only when all dependencies (DB, Redis, cache warming) are healthy prevents HTTP 502/503 errors during deployment.

---

## Evidence 4

Claim: Graceful termination involves sending a termination signal (SIGTERM), waiting for a configured grace period to let active requests complete, and detaching the endpoint from the router/load balancer before sending SIGKILL.

Evidence: "Typically, with this graceful termination of the pod, kubelet makes requests to the container runtime to attempt to stop the containers in the pod by first sending a TERM (aka. SIGTERM) signal, with a grace period timeout, to the main process in each container... Once the grace period has expired, the KILL signal is sent to any remaining processes, and the Pod is then deleted from the API Server."

Source: Pod Lifecycle
URL: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/
Confidence: HIGH
Corroborated By: General POSIX process management standards and standard container orchestrators
Notes: Enables connection draining so in-flight requests (such as checkouts or file uploads) do not get dropped.

---

## Evidence 5

Claim: Background queue workers are long-lived processes that hold old code in memory and must be gracefully restarted after deployment to pick up code changes without losing or aborting running jobs.

Evidence: "Since queue workers are long-lived processes, they will not notice changes to your code without being restarted. So, the simplest way to deploy an application using queue workers is to restart the workers during your deployment process. You may gracefully restart all of the workers by issuing the `queue:restart` command... This command will instruct all queue workers to gracefully exit after they finish processing their current job so that no existing jobs are lost."

Source: Queues - Queue Workers and Deployment
URL: https://laravel.com/docs/11.x/queues
Confidence: HIGH
Corroborated By: Supervisor process monitor documentation / standard daemon management patterns
Notes: Requires a process manager (like Supervisor or Kubernetes Pod lifecycle management) to re-spawn the worker process with the new version.