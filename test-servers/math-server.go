package main

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Math Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add calculate tool
	calculateTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform: add, subtract, multiply, divide"),
		),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	// Add square root tool
	sqrtTool := mcp.NewTool("sqrt",
		mcp.WithDescription("Calculate square root of a number"),
		mcp.WithNumber("number",
			mcp.Required(),
			mcp.Description("Number to calculate square root of"),
		),
	)

	// Add power tool
	powerTool := mcp.NewTool("power",
		mcp.WithDescription("Calculate a number raised to a power"),
		mcp.WithNumber("base",
			mcp.Required(),
			mcp.Description("Base number"),
		),
		mcp.WithNumber("exponent",
			mcp.Required(),
			mcp.Description("Exponent"),
		),
	)

	// Add tool handlers
	s.AddTool(calculateTool, calculateHandler)
	s.AddTool(sqrtTool, sqrtHandler)
	s.AddTool(powerTool, powerHandler)

	// Start stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Math server error: %v\n", err)
		os.Exit(1)
	}
}

func calculateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	operation, err := request.RequireString("operation")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	a, err := request.RequireFloat("a")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	b, err := request.RequireFloat("b")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var result float64
	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return mcp.NewToolResultError("Division by zero"), nil
		}
		result = a / b
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s", operation)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
}

func sqrtHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	number, err := request.RequireFloat("number")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if number < 0 {
		return mcp.NewToolResultError("Cannot calculate square root of negative number"), nil
	}

	result := math.Sqrt(number)
	return mcp.NewToolResultText(fmt.Sprintf("%.6f", result)), nil
}

func powerHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	base, err := request.RequireFloat("base")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	exponent, err := request.RequireFloat("exponent")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := math.Pow(base, exponent)
	return mcp.NewToolResultText(fmt.Sprintf("%.6f", result)), nil
}