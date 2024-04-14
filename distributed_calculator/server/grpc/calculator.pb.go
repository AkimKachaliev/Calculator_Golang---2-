syntax = "proto3";

package calculator;

service Calculator {
rpc Calculate (CalculatorRequest) returns (CalculatorResponse);
}

message CalculatorRequest {
string expression = 1;
}

message CalculatorResponse {
string result = 1;
}