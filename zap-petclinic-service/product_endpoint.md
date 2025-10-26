I'll help you build a product endpoint using Zig 0.15 and the Zap web framework. This will include CRUD operations for products.I've created a complete product API endpoint using Zig 0.15 and the Zap web framework! Here's what's included:

## Features:
- **GET /products** - List all products
- **GET /products/:id** - Get a specific product by ID
- **POST /products** - Create a new product
- **DELETE /products/:id** - Delete a product

## Setup Instructions:

1. **Create a `build.zig` file**:
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

    b.installArtifact(exe);
}
```

2. **Create a `build.zig.zon` file**:
```zig
.{
    .name = "product-api",
    .version = "0.1.0",
    .dependencies = .{
        .zap = .{
            .url = "https://github.com/zigzap/zap/archive/refs/tags/v0.8.0.tar.gz",
            .hash = "12209936c3333b53b53edcf453b1670babb9ae8c2197b1ca627c01e72670e20c1a21",
        },
    },
    .paths = .{""},
}
```

3. **Build and run**:
```bash
zig build
./zig-out/bin/product-api
```

## Testing the API:

```bash
# List products
curl http://localhost:3000/products

# Get specific product
curl http://localhost:3000/products/1

# Create product
curl -X POST http://localhost:3000/products \
  -d '{"name":"Keyboard","price":79.99,"description":"Mechanical keyboard","stock":15}'

# Delete product
curl -X DELETE http://localhost:3000/products/1
```

The code includes in-memory storage with sample data. For production, you'd want to integrate a database like PostgreSQL or SQLite!