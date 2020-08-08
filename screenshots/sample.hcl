# Map size rows x columns.
rows = 12
cols = 7

# Map margin 
margin = 20

# Map background
background = "#fafafa"

# Should draw the grid? (useful for debug)
grid = false

# Should draw borders? (useful for debug)
border = false 

# Should show the grid coordinates? (useful for debug)
hints = false

# Draw a rectangular frame
tile "frame" "" {
    left = "1"
    top = "0"
    right = "11"
    bottom = "6"
    color = "#b9b9b099"
    stroke = true
    stroke_width = 2
    dashes = 8
}

# AWS API Gateway 
# use embedded icon library.
tile "icon" "agw" {
    row = 9
    col = 4
    uri = "assets://aws_api_gateway"
}

# AWS Lambda Authorizer
tile "icon" "lambda1" {
    row = "${subtract(row("agw"), 1)}"
    col = "${subtract(col("agw"), 2)}"
    uri = "assets://aws_lambda"
}

# AWS Lambda function
tile "icon" "lambda2" {
    row = "${add(row("agw"), 1)}"
    col = "${subtract(col("agw"), 2)}"
    uri = "assets://aws_lambda"
}

# Connection
tile "tee_right" "" {
    row = "${row("agw")}"
    col = "${subtract(col("agw"), 1)}"
}

# Connection
tile "elbow_left_down" "" {
    row = "${row("lambda1")}"
    col = "${subtract(col("agw"), 1)}"
    arrow_left = true
}

# Connection
tile "elbow_left_up" "" {
    row = "${row("lambda2")}"
    col = "${subtract(col("agw"), 1)}"
    arrow_left = true
}

# AWS ECS
tile "icon" "ecs" {
    row = "${subtract(row("agw"), 5)}"
    col = "${col("agw")}"
    uri = "assets://aws_elastic_container_service"
}

# AWS RDS MySQL 
tile "icon" "mysql" {
    row = "${subtract(row("ecs"), 2)}"
    col = "${subtract(col("ecs"), 2)}"
    uri = "assets://aws_rds_mysql_instance"
}

# AWS Simple Storage Service (S3) 
tile "icon" "s3" {
    row = "${row("ecs")}"
    col = "${col("mysql")}"
    uri = "assets://aws_simple_storage_service_s3"
}

# AWS elasticache for Redis
tile "icon" "redis" {
    row = "${add(row("ecs"), 2)}"
    col = "${col("s3")}"
    uri = "assets://aws_elasticache_for_redis"
}

# Connection
tile "cross" "" {
    row = "${row("ecs")}"
    col = "${subtract(col("ecs"), 1)}"
    arrow_left = true
}

# Connection
tile "vertical_line" "" {
    row = "${subtract(row("ecs"), 1)}"
    col = "${subtract(col("ecs"), 1)}"
}

# Connection
tile "vertical_line" "" {
    row = "${add(row("ecs"), 1)}"
    col = "${subtract(col("ecs"), 1)}"
}

# Connection
tile "elbow_left_down" "" {
    row = "${subtract(row("ecs"), 2)}"
    col = "${subtract(col("ecs"), 1)}"
    arrow_left = true
}

# Connection
tile "elbow_left_up" "" {
    row = "${add(row("ecs"), 2)}"
    col = "${subtract(col("ecs"), 1)}"
    arrow_left = true
}

# Connection
tile "vertical_line" "" {
    row = "${subtract(row("agw"), 1)}"
    col = "${col("agw")}"
}

# Connection 
tile "vertical_line" "" {
    row = "${subtract(row("agw"), 2)}"
    col = "${col("agw")}"
}

# Connection
tile "vertical_line" "" {
    row = "${subtract(row("agw"), 3)}"
    col = "${col("agw")}"
}

# Connection
tile "vertical_line" "" {
    row = "${subtract(row("agw"), 4)}"
    col = "${col("agw")}"
    arrow_up = true
}

# Label 
tile "label" "" {
    row = "${subtract(row("agw"), 8)}"
    col = "${subtract(col("agw"), 3)}"
    text = "AWS"
    font_size = 15
    background = "#e5d512"
}

# Label
tile "label" "" {
    row = "${add(row("agw"), 2)}"
    col = "${subtract(col("agw"), 3)}"
    text = "token manager account"
    font_size = 14
    color = "#ffffff"
    background = "#68675caa"
}
