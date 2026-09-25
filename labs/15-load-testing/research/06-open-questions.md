# Open Questions

- What are the SLA limits enforced by third-party dependencies (e.g., Payment Gateways, WhatsApp API) in the Booking Bengkel architecture?
- Does the system employ circuit breakers for the WhatsApp integration to prevent thread exhaustion during external API timeouts?
- Is there a read-replica database available to offload analytical or non-critical reads during peak transaction processing?
- What are the current production hardware specs, and what is the expected maximum concurrent user count during peak hours?
- Are auto-scaling configurations implemented for application servers under high CPU/memory utilization?
