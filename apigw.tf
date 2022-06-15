resource "aws_apigatewayv2_api" "blossom" {
  name = "${local.prefix}-gw"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "blossom" {
  name = "${local.prefix}-gw-stage"

  api_id = aws_apigatewayv2_api.blossom
  auto_deploy = true
}
