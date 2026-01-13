package com.example;

import com.google.adk.tools.BaseTool;
import com.google.adk.tools.FunctionTool;
import java.lang.reflect.Method;

public class CheckApi {
    public static void main(String[] args) {
        System.out.println("=== BaseTool methods ===");
        for (Method m : BaseTool.class.getMethods()) {
            System.out.println("  " + m.getName() + "() -> " + m.getReturnType().getSimpleName());
        }
        
        System.out.println("\n=== FunctionTool methods ===");
        for (Method m : FunctionTool.class.getMethods()) {
            if (m.getDeclaringClass() == FunctionTool.class) {
                System.out.println("  " + m.getName() + "() -> " + m.getReturnType().getSimpleName());
            }
        }
        
        // Try to create a FunctionTool and inspect it
        System.out.println("\n=== Testing FunctionTool.create ===");
        try {
            BaseTool tool = FunctionTool.create(CheckApi.class, "testMethod");
            System.out.println("Tool created: " + tool.getClass().getName());
            for (Method m : tool.getClass().getMethods()) {
                if (m.getName().toLowerCase().contains("function") || m.getName().toLowerCase().contains("declaration")) {
                    System.out.println("  Found: " + m.getName() + "() -> " + m.getReturnType());
                }
            }
        } catch (Exception e) {
            System.out.println("Error: " + e.getMessage());
        }
        
        System.exit(0);
    }
    
    @com.google.adk.tools.Annotations.Schema(description = "Test method")
    public static String testMethod(@com.google.adk.tools.Annotations.Schema(description = "input") String input) {
        return "test: " + input;
    }
}
