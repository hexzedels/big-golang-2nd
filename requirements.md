Functional Requirements:
    1. Users can submit one-time or periodic jobs for execution.
    2. Users can cancel the submitted jobs.
    3. The system should distribute jobs across multiple worker nodes for execution.
    4. The system should provide monitoring of job status (queued, running, completed, failed).
    5. The system should prevent the same job from being executed multiple times concurrently.
    6. The system should reschedule job in case of worker failure.

Non-Functional Requirements:
    1. Scalability: The system should be able to schedule and execute 10k/100k/1kk of jobs.
    2. High Availability: The system should be fault-tolerant with no single point of failure. If a worker node fails, the system should reschedule the job to other available nodes.
    3. Latency: Jobs should be scheduled and executed with minimal delay. (1-3s)
    4. Consistency: Job results should be consistent, ensuring that jobs are executed once (or with minimal duplication).