// Project Structure:
// product-api/
// ├── build.zig
// ├── build.zig.zon
// ├── src/
// │   ├── main.zig
// │   ├── config/
// │   │   └── config.zig
// │   ├── models/
// │   │   └── product.zig
// │   ├── repositories/
// │   │   └── product_repository.zig
// │   ├── handlers/
// │   │   └── product_handler.zig
// │   └── database/
// │       └── sqlite.zig
// ├── config.json
// └── README.md

// ============================================================================
// File: build.zig
// ============================================================================
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

    // Add zap dependency
    const zap = b.dependency("zap", .{
        .target = target,
        .optimize = optimize,
    });
    exe.root_module.addImport("zap", zap.module("zap"));

    // Link SQLite
    exe.linkLibC();
    exe.linkSystemLibrary("sqlite3");

    b.installArtifact(exe);

    const run_cmd = b.addRunArtifact(exe);
    run_cmd.step.dependOn(b.getInstallStep());

    if (b.args) |args| {
        run_cmd.addArgs(args);
    }

    const run_step = b.step("run", "Run the app");
    run_step.dependOn(&run_cmd.step);
}

// ============================================================================
// File: build.zig.zon
// ============================================================================
// .{
//     .name = "product-api",
//     .version = "0.1.0",
//     .dependencies = .{
//         .zap = .{
//             .url = "https://github.com/zigzap/zap/archive/refs/tags/v0.8.0.tar.gz",
//             .hash = "12209936c3333b53b53edcf453b1670babb9ae8c2197b1ca627c01e72670e20c1a21",
//         },
//     },
//     .paths = .{
//         "build.zig",
//         "build.zig.zon",
//         "src",
//     },
// }

// ============================================================================
// File: src/config/config.zig
// ============================================================================
const std = @import("std");

pub const Config = struct {
    server: ServerConfig,
    database: DatabaseConfig,

    pub const ServerConfig = struct {
        host: []const u8,
        port: u16,
    };

    pub const DatabaseConfig = struct {
        path: []const u8,
    };

    pub fn load(allocator: std.mem.Allocator, path: []const u8) !Config {
        const file = try std.fs.cwd().openFile(path, .{});
        defer file.close();

        const content = try file.readToEndAlloc(allocator, 1024 * 1024);
        defer allocator.free(content);

        const parsed = try std.json.parseFromSlice(Config, allocator, content, .{});
        defer parsed.deinit();

        // Duplicate strings for owned config
        const server_host = try allocator.dupe(u8, parsed.value.server.host);
        const db_path = try allocator.dupe(u8, parsed.value.database.path);

        return Config{
            .server = .{
                .host = server_host,
                .port = parsed.value.server.port,
            },
            .database = .{
                .path = db_path,
            },
        };
    }

    pub fn deinit(self: *Config, allocator: std.mem.Allocator) void {
        allocator.free(self.server.host);
        allocator.free(self.database.path);
    }
};

// ============================================================================
// File: src/models/product.zig
// ============================================================================
const std = @import("std");

pub const Product = struct {
    id: i64,
    name: []const u8,
    description: []const u8,
    price: f64,
    stock: i32,
    created_at: []const u8,

    pub fn deinit(self: *Product, allocator: std.mem.Allocator) void {
        allocator.free(self.name);
        allocator.free(self.description);
        allocator.free(self.created_at);
    }
};

pub const CreateProductRequest = struct {
    name: []const u8,
    description: []const u8,
    price: f64,
    stock: i32,
};

pub const UpdateProductRequest = struct {
    name: ?[]const u8 = null,
    description: ?[]const u8 = null,
    price: ?f64 = null,
    stock: ?i32 = null,
};

// ============================================================================
// File: src/database/sqlite.zig
// ============================================================================
const std = @import("std");
const c = @cImport({
    @cInclude("sqlite3.h");
});

pub const Database = struct {
    db: *c.sqlite3,
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator, path: []const u8) !Database {
        var db: ?*c.sqlite3 = null;
        const path_z = try allocator.dupeZ(u8, path);
        defer allocator.free(path_z);

        const result = c.sqlite3_open(path_z.ptr, &db);
        if (result != c.SQLITE_OK) {
            std.debug.print("Failed to open database: {s}\n", .{c.sqlite3_errmsg(db)});
            return error.DatabaseOpenFailed;
        }

        return Database{
            .db = db.?,
            .allocator = allocator,
        };
    }

    pub fn deinit(self: *Database) void {
        _ = c.sqlite3_close(self.db);
    }

    pub fn exec(self: *Database, sql: []const u8) !void {
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var err_msg: ?[*:0]u8 = null;
        const result = c.sqlite3_exec(self.db, sql_z.ptr, null, null, &err_msg);

        if (result != c.SQLITE_OK) {
            if (err_msg) |msg| {
                std.debug.print("SQL error: {s}\n", .{msg});
                c.sqlite3_free(msg);
            }
            return error.SqlExecFailed;
        }
    }

    pub fn createTables(self: *Database) !void {
        const create_products_table =
            \\CREATE TABLE IF NOT EXISTS products (
            \\    id INTEGER PRIMARY KEY AUTOINCREMENT,
            \\    name TEXT NOT NULL,
            \\    description TEXT,
            \\    price REAL NOT NULL,
            \\    stock INTEGER NOT NULL DEFAULT 0,
            \\    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            \\);
        ;
        try self.exec(create_products_table);
        std.debug.print("Database tables created successfully\n", .{});
    }
};

// ============================================================================
// File: src/repositories/product_repository.zig
// ============================================================================
const std = @import("std");
const Product = @import("../models/product.zig").Product;
const CreateProductRequest = @import("../models/product.zig").CreateProductRequest;
const UpdateProductRequest = @import("../models/product.zig").UpdateProductRequest;
const Database = @import("../database/sqlite.zig").Database;
const c = @cImport({
    @cInclude("sqlite3.h");
});

pub const ProductRepository = struct {
    db: *Database,
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator, db: *Database) ProductRepository {
        return ProductRepository{
            .db = db,
            .allocator = allocator,
        };
    }

    pub fn create(self: *ProductRepository, req: CreateProductRequest) !i64 {
        const sql = "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)";
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var stmt: ?*c.sqlite3_stmt = null;
        var result = c.sqlite3_prepare_v2(self.db.db, sql_z.ptr, -1, &stmt, null);

        if (result != c.SQLITE_OK) {
            return error.PrepareFailed;
        }
        defer _ = c.sqlite3_finalize(stmt);

        const name_z = try self.allocator.dupeZ(u8, req.name);
        defer self.allocator.free(name_z);
        const desc_z = try self.allocator.dupeZ(u8, req.description);
        defer self.allocator.free(desc_z);

        _ = c.sqlite3_bind_text(stmt, 1, name_z.ptr, -1, c.SQLITE_TRANSIENT);
        _ = c.sqlite3_bind_text(stmt, 2, desc_z.ptr, -1, c.SQLITE_TRANSIENT);
        _ = c.sqlite3_bind_double(stmt, 3, req.price);
        _ = c.sqlite3_bind_int(stmt, 4, req.stock);

        result = c.sqlite3_step(stmt);
        if (result != c.SQLITE_DONE) {
            return error.InsertFailed;
        }

        return c.sqlite3_last_insert_rowid(self.db.db);
    }

    pub fn findAll(self: *ProductRepository) !std.ArrayList(Product) {
        const sql = "SELECT id, name, description, price, stock, created_at FROM products";
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var stmt: ?*c.sqlite3_stmt = null;
        var result = c.sqlite3_prepare_v2(self.db.db, sql_z.ptr, -1, &stmt, null);

        if (result != c.SQLITE_OK) {
            return error.PrepareFailed;
        }
        defer _ = c.sqlite3_finalize(stmt);

        var products = std.ArrayList(Product).init(self.allocator);

        while (c.sqlite3_step(stmt) == c.SQLITE_ROW) {
            const id = c.sqlite3_column_int64(stmt, 0);
            const name_ptr = c.sqlite3_column_text(stmt, 1);
            const desc_ptr = c.sqlite3_column_text(stmt, 2);
            const price = c.sqlite3_column_double(stmt, 3);
            const stock = c.sqlite3_column_int(stmt, 4);
            const created_ptr = c.sqlite3_column_text(stmt, 5);

            const name = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(name_ptr))));
            const desc = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(desc_ptr))));
            const created = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(created_ptr))));

            try products.append(Product{
                .id = id,
                .name = name,
                .description = desc,
                .price = price,
                .stock = stock,
                .created_at = created,
            });
        }

        return products;
    }

    pub fn findById(self: *ProductRepository, id: i64) !?Product {
        const sql = "SELECT id, name, description, price, stock, created_at FROM products WHERE id = ?";
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var stmt: ?*c.sqlite3_stmt = null;
        var result = c.sqlite3_prepare_v2(self.db.db, sql_z.ptr, -1, &stmt, null);

        if (result != c.SQLITE_OK) {
            return error.PrepareFailed;
        }
        defer _ = c.sqlite3_finalize(stmt);

        _ = c.sqlite3_bind_int64(stmt, 1, id);

        if (c.sqlite3_step(stmt) == c.SQLITE_ROW) {
            const name_ptr = c.sqlite3_column_text(stmt, 1);
            const desc_ptr = c.sqlite3_column_text(stmt, 2);
            const price = c.sqlite3_column_double(stmt, 3);
            const stock = c.sqlite3_column_int(stmt, 4);
            const created_ptr = c.sqlite3_column_text(stmt, 5);

            const name = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(name_ptr))));
            const desc = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(desc_ptr))));
            const created = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(created_ptr))));

            return Product{
                .id = id,
                .name = name,
                .description = desc,
                .price = price,
                .stock = stock,
                .created_at = created,
            };
        }

        return null;
    }

    pub fn deleteById(self: *ProductRepository, id: i64) !bool {
        const sql = "DELETE FROM products WHERE id = ?";
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var stmt: ?*c.sqlite3_stmt = null;
        var result = c.sqlite3_prepare_v2(self.db.db, sql_z.ptr, -1, &stmt, null);

        if (result != c.SQLITE_OK) {
            return error.PrepareFailed;
        }
        defer _ = c.sqlite3_finalize(stmt);

        _ = c.sqlite3_bind_int64(stmt, 1, id);
        result = c.sqlite3_step(stmt);

        if (result != c.SQLITE_DONE) {
            return error.DeleteFailed;
        }

        return c.sqlite3_changes(self.db.db) > 0;
    }
};

// ============================================================================
// File: src/handlers/product_handler.zig
// ============================================================================
const std = @import("std");
const zap = @import("zap");
const Product = @import("../models/product.zig").Product;
const CreateProductRequest = @import("../models/product.zig").CreateProductRequest;
const ProductRepository = @import("../repositories/product_repository.zig").ProductRepository;

pub const ProductHandler = struct {
    repo: *ProductRepository,
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator, repo: *ProductRepository) ProductHandler {
        return ProductHandler{
            .repo = repo,
            .allocator = allocator,
        };
    }

    pub fn getAllProducts(self: *ProductHandler, r: zap.Request) void {
        var products = self.repo.findAll() catch {
            r.setStatus(.internal_server_error);
            r.sendBody("{\"error\":\"Failed to fetch products\"}") catch return;
            return;
        };
        defer {
            for (products.items) |*p| {
                p.deinit(self.allocator);
            }
            products.deinit();
        }

        var response = std.ArrayList(u8).init(self.allocator);
        defer response.deinit();

        response.appendSlice("{\"products\":[") catch return;
        for (products.items, 0..) |product, i| {
            if (i > 0) response.appendSlice(",") catch return;
            std.fmt.format(response.writer(), "{{\"id\":{d},\"name\":\"{s}\",\"description\":\"{s}\",\"price\":{d:.2},\"stock\":{d},\"created_at\":\"{s}\"}}", .{
                product.id,
                product.name,
                product.description,
                product.price,
                product.stock,
                product.created_at,
            }) catch return;
        }
        response.appendSlice("]}") catch return;

        r.setStatus(.ok);
        r.setContentType(.JSON) catch return;
        r.sendBody(response.items) catch return;
    }

    pub fn getProductById(self: *ProductHandler, r: zap.Request) void {
        const path = r.path orelse {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid path\"}") catch return;
            return;
        };

        var it = std.mem.splitSequence(u8, path, "/");
        var id_str: ?[]const u8 = null;
        var count: usize = 0;
        while (it.next()) |segment| {
            count += 1;
            if (count == 3) id_str = segment;
        }

        const id = std.fmt.parseInt(i64, id_str orelse "", 10) catch {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid product ID\"}") catch return;
            return;
        };

        var product = self.repo.findById(id) catch {
            r.setStatus(.internal_server_error);
            r.sendBody("{\"error\":\"Failed to fetch product\"}") catch return;
            return;
        };

        if (product) |*p| {
            defer p.deinit(self.allocator);

            var response = std.ArrayList(u8).init(self.allocator);
            defer response.deinit();

            std.fmt.format(response.writer(), "{{\"id\":{d},\"name\":\"{s}\",\"description\":\"{s}\",\"price\":{d:.2},\"stock\":{d},\"created_at\":\"{s}\"}}", .{
                p.id,
                p.name,
                p.description,
                p.price,
                p.stock,
                p.created_at,
            }) catch return;

            r.setStatus(.ok);
            r.setContentType(.JSON) catch return;
            r.sendBody(response.items) catch return;
        } else {
            r.setStatus(.not_found);
            r.sendBody("{\"error\":\"Product not found\"}") catch return;
        }
    }

    pub fn createProduct(self: *ProductHandler, r: zap.Request) void {
        const body = r.body orelse {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Request body required\"}") catch return;
            return;
        };

        const parsed = std.json.parseFromSlice(CreateProductRequest, self.allocator, body, .{}) catch {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid JSON\"}") catch return;
            return;
        };
        defer parsed.deinit();

        const id = self.repo.create(parsed.value) catch {
            r.setStatus(.internal_server_error);
            r.sendBody("{\"error\":\"Failed to create product\"}") catch return;
            return;
        };

        var response = std.ArrayList(u8).init(self.allocator);
        defer response.deinit();

        std.fmt.format(response.writer(), "{{\"id\":{d},\"message\":\"Product created successfully\"}}", .{id}) catch return;

        r.setStatus(.created);
        r.setContentType(.JSON) catch return;
        r.sendBody(response.items) catch return;
    }

    pub fn deleteProduct(self: *ProductHandler, r: zap.Request) void {
        const path = r.path orelse {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid path\"}") catch return;
            return;
        };

        var it = std.mem.splitSequence(u8, path, "/");
        var id_str: ?[]const u8 = null;
        var count: usize = 0;
        while (it.next()) |segment| {
            count += 1;
            if (count == 3) id_str = segment;
        }

        const id = std.fmt.parseInt(i64, id_str orelse "", 10) catch {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid product ID\"}") catch return;
            return;
        };

        const deleted = self.repo.deleteById(id) catch {
            r.setStatus(.internal_server_error);
            r.sendBody("{\"error\":\"Failed to delete product\"}") catch return;
            return;
        };

        if (deleted) {
            r.setStatus(.ok);
            r.sendBody("{\"message\":\"Product deleted successfully\"}") catch return;
        } else {
            r.setStatus(.not_found);
            r.sendBody("{\"error\":\"Product not found\"}") catch return;
        }
    }
};

// ============================================================================
// File: src/main.zig
// ============================================================================
const std = @import("std");
const zap = @import("zap");
const Config = @import("config/config.zig").Config;
const Database = @import("database/sqlite.zig").Database;
const ProductRepository = @import("repositories/product_repository.zig").ProductRepository;
const ProductHandler = @import("handlers/product_handler.zig").ProductHandler;

var gpa = std.heap.GeneralPurposeAllocator(.{}){};
var db: Database = undefined;
var repo: ProductRepository = undefined;
var handler: ProductHandler = undefined;

fn on_request(r: zap.Request) void {
    const path = r.path orelse {
        r.setStatus(.bad_request);
        r.sendBody("{\"error\":\"Invalid path\"}") catch return;
        return;
    };

    if (std.mem.eql(u8, path, "/api/products")) {
        if (r.method) |method| {
            if (method == .GET) {
                handler.getAllProducts(r);
                return;
            } else if (method == .POST) {
                handler.createProduct(r);
                return;
            }
        }
    } else if (std.mem.startsWith(u8, path, "/api/products/")) {
        if (r.method) |method| {
            if (method == .GET) {
                handler.getProductById(r);
                return;
            } else if (method == .DELETE) {
                handler.deleteProduct(r);
                return;
            }
        }
    } else if (std.mem.eql(u8, path, "/health")) {
        r.setStatus(.ok);
        r.sendBody("{\"status\":\"ok\"}") catch return;
        return;
    }

    r.setStatus(.not_found);
    r.sendBody("{\"error\":\"Not found\"}") catch return;
}

pub fn main() !void {
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    // Load configuration
    var config = try Config.load(allocator, "config.json");
    defer config.deinit(allocator);

    std.debug.print("Configuration loaded:\n", .{});
    std.debug.print("  Server: {s}:{d}\n", .{ config.server.host, config.server.port });
    std.debug.print("  Database: {s}\n", .{config.database.path});

    // Initialize database
    db = try Database.init(allocator, config.database.path);
    defer db.deinit();
    try db.createTables();

    // Initialize repository, handler and dependency injection
    repo = ProductRepository.init(allocator, &db);
    handler = ProductHandler.init(allocator, &repo);

    // Setup Zap server
    var listener = zap.HttpListener.init(.{
        .port = config.server.port,
        .on_request = on_request,
        .log = true,
    });
    try listener.listen();

    std.debug.print("\n🚀 Server running on http://{s}:{d}\n", .{ config.server.host, config.server.port });
    std.debug.print("\nAvailable endpoints:\n", .{});
    std.debug.print("  GET    /health              - Health check\n", .{});
    std.debug.print("  GET    /api/products        - Get all products\n", .{});
    std.debug.print("  GET    /api/products/:id    - Get product by ID\n", .{});
    std.debug.print("  POST   /api/products        - Create new product\n", .{});
    std.debug.print("  DELETE /api/products/:id    - Delete product\n", .{});

    zap.start(.{
        .threads = 4,
        .workers = 2,
    });
}

// ============================================================================
// File: config.json
// ============================================================================
// {
//   "server": {
//     "host": "127.0.0.1",
//     "port": 3000
//   },
//   "database": {
//     "path": "./products.db"
//   }
// }

// ============================================================================
// File: README.md
// ============================================================================
// # Zig Product API
//
// A RESTful API built with Zig, Zap web framework, and SQLite.
//
// ## Prerequisites
//
// - Zig 0.13.0 or later
// - SQLite3 development libraries
//
// ## Installation
//
// ### Install SQLite (Ubuntu/Debian)
// ```bash
// sudo apt-get install libsqlite3-dev
// ```
//
// ### Install SQLite (macOS)
// ```bash
// brew install sqlite3
// ```
//
// ## Build & Run
//
// ```bash
// # Build the project
// zig build
//
// # Run the server
// zig build run
// ```
//
// ## API Endpoints
//
// ### Health Check
// ```bash
// curl http://localhost:3000/health
// ```
//
// ### Get All Products
// ```bash
// curl http://localhost:3000/api/products
// ```
//
// ### Get Product by ID
// ```bash
// curl http://localhost:3000/api/products/1
// ```
//
// ### Create Product
// ```bash
// curl -X POST http://localhost:3000/api/products \
//   -H "Content-Type: application/json" \
//   -d '{
//     "name": "Laptop",
//     "description": "High-performance laptop",
//     "price": 999.99,
//     "stock": 50
//   }'
// ```
//
// ### Delete Product
// ```bash
// curl -X DELETE http://localhost:3000/api/products/1
// ```
//
// ## Project Structure
//
// ```
// product-api/
// ├── build.zig              # Build configuration
// ├── build.zig.zon          # Dependencies
// ├── config.json            # Application configuration
// ├── src/
// │   ├── main.zig           # Application entry point
// │   ├── config/
// │   │   └── config.zig     # Configuration management
// │   ├── models/
// │   │   └── product.zig    # Product data models
// │   ├── repositories/
// │   │   └── product_repository.zig  # Database operations
// │   ├── handlers/
// │   │   └── product_handler.zig     # HTTP request handlers
// │   └── database/
// │       └── sqlite.zig     # SQLite wrapper
// └── README.md
// ```