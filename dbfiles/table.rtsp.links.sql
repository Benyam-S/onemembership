CREATE TABLE rtsp_links (
    resource_id VARCHAR(255) NOT NULL,
    plan_id VARCHAR(255) NOT NULL,
    UNIQUE KEY unique_resource_to_subscription_plan_relation (resource_id, plan_id),
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (plan_id) REFERENCES subscription_plans(id) ON DELETE CASCADE ON UPDATE CASCADE
);