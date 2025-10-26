I'll help you build a product endpoint using Zig 0.15, Zap web framework, and SQLite. This will include CRUD operations for products.I've created a complete Product API using Zig 0.15, Zap web framework, and SQLite! Here's what's included:

## Features

**Full CRUD Operations:**
- `GET /api/products` - Get all products
- `GET /api/products/:id` - Get a single product by ID
- `POST /api/products` - Create a new product
- `PUT /api/products/:id` - Update an existing product
- `DELETE /api/products/:id` - Delete a product

**Product Schema:**
- `id` - Auto-incrementing primary key
- `name` - Product name (required)
- `description` - Product description
- `price` - Product price (required)
- `stock` - Stock quantity (required)

## Setup

You'll need to create a `build.zig` file to include the dependencies:

```zig
const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const exe = b.addExecutable(.{
        .name = "product-api",
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    const zap = b.dependency("zap", .{
        .target = target,
        .optimize = optimize,
    });
    exe.root_module.addImport("zap", zap.module("zap"));

    const sqlite = b.dependency("sqlite", .{
        .target = target,
        .optimize = optimize,
    });
    exe.root_module.addImport("sqlite", sqlite.module("sqlite"));

    b.installArtifact(exe);
}
```

## Example Usage

**Create a product:**
```bash
curl -X POST http://localhost:3000/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"Gaming laptop","price":1299.99,"stock":10}'
```

**Get all products:**
```bash
curl http://localhost:3000/api/products
```

The server will create a `products.db` SQLite database file automatically on first run!