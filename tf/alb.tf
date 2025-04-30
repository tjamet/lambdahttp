resource "aws_lb" "rest-api-alb" {
  name               = "rest-api-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = data.aws_subnets.default.ids
}

resource "aws_security_group" "alb_sg" {
  name        = "rest-api-alb-sg"
  description = "Security group for ALB"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb_target_group" "rest-api-tg" {
  name        = "rest-api-tg"
  target_type = "lambda"
}

resource "aws_lb_target_group_attachment" "rest-api-tg-attachment" {
  target_group_arn = aws_lb_target_group.rest-api-tg.arn
  target_id        = aws_lambda_function.rest-api.arn
}

resource "aws_lambda_permission" "allow_alb" {
  statement_id  = "AllowALBInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.rest-api.function_name
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = aws_lb_target_group.rest-api-tg.arn
}

resource "aws_lb_listener" "rest-api-listener" {
  load_balancer_arn = aws_lb.rest-api-alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.rest-api-tg.arn
  }
} 