const std = @import("std");
const zap = @import("zap");

// Product structure
const Product = struct {
    id: u32,
    name: []const u8,
    price: f64,
    description: []const u8,
    stock: u32,
};

// In-memory storage (in production, use a database)
var products = std.ArrayListUnmanaged(Product){};
var next_id: u32 = 1;
var gpa = std.heap.GeneralPurposeAllocator(.{}){};

// Handler for GET /products - List all products
fn getProducts(r: zap.Request) void {
    r.setContentType(.JSON) catch return;

    var buf = std.ArrayList(u8).init(gpa.allocator());
    defer buf.deinit();

    var writer = buf.writer();
    writer.writeAll("[") catch return;

    for (products.items, 0..) |p, i| {
        if (i > 0) writer.writeAll(",") catch return;
        std.fmt.format(writer,
            \\{{"id":{d},"name":"{s}","price":{d:.2},"description":"{s}","stock":{d}}}
        , .{p.id, p.name, p.price, p.description, p.stock}) catch return;
    }

    writer.writeAll("]") catch return;
    r.sendBody(buf.items) catch return;
}

// Handler for GET /products/:id - Get single product
fn getProduct(r: zap.Request) void {
    if (r.getParamStr("id")) |id_str| {
        const id = std.fmt.parseInt(u32, id_str, 10) catch {
            r.setStatus(.bad_request) catch return;
            r.sendBody("Invalid ID") catch return;
            return;
        };

        for (products.items) |p| {
            if (p.id == id) {
                r.setContentType(.JSON) catch return;
                var buf: [512]u8 = undefined;
                const json = std.fmt.bufPrint(&buf,
                    \\{{"id":{d},"name":"{s}","price":{d:.2},"description":"{s}","stock":{d}}}
                , .{p.id, p.name, p.price, p.description, p.stock}) catch return;
                r.sendBody(json) catch return;
                return;
            }
        }

        r.setStatus(.not_found) catch return;
        r.sendBody("Product not found") catch return;
    }
}

// Handler for POST /products - Create product
fn createProduct(r: zap.Request) void {
    if (r.body) |body| {
        // Simple JSON parsing (in production, use a proper JSON library)
        // Expected format: {"name":"Product","price":99.99,"description":"Desc","stock":10}

        var name: []const u8 = "";
        var price: f64 = 0;
        var description: []const u8 = "";
        var stock: u32 = 0;

        // Basic parsing (this is simplified - use std.json in real code)
        var iter = std.mem.split(u8, body, ",");
        while (iter.next()) |field| {
            if (std.mem.indexOf(u8, field, "\"name\"") != null) {
                if (std.mem.lastIndexOf(u8, field, "\"")) |end| {
                    if (std.mem.indexOf(u8, field[0..end], ":\"")) |start| {
                        name = field[start+2..end];
                    }
                }
            } else if (std.mem.indexOf(u8, field, "\"price\"") != null) {
                if (std.mem.indexOf(u8, field, ":")) |colon| {
                    const val = std.mem.trim(u8, field[colon+1..], " \t\r\n}");
                    price = std.fmt.parseFloat(f64, val) catch 0;
                }
            } else if (std.mem.indexOf(u8, field, "\"description\"") != null) {
                if (std.mem.lastIndexOf(u8, field, "\"")) |end| {
                    if (std.mem.indexOf(u8, field[0..end], ":\"")) |start| {
                        description = field[start+2..end];
                    }
                }
            } else if (std.mem.indexOf(u8, field, "\"stock\"") != null) {
                if (std.mem.indexOf(u8, field, ":")) |colon| {
                    const val = std.mem.trim(u8, field[colon+1..], " \t\r\n}");
                    stock = std.fmt.parseInt(u32, val, 10) catch 0;
                }
            }
        }

        const product = Product{
            .id = next_id,
            .name = name,
            .price = price,
            .description = description,
            .stock = stock,
        };

        next_id += 1;
        products.append(product) catch return;

        r.setStatus(.created) catch return;
        r.setContentType(.JSON) catch return;
        var buf: [512]u8 = undefined;
        const json = std.fmt.bufPrint(&buf,
            \\{{"id":{d},"name":"{s}","price":{d:.2},"description":"{s}","stock":{d}}}
        , .{product.id, product.name, product.price, product.description, product.stock}) catch return;
        r.sendBody(json) catch return;
    }
}

// Handler for DELETE /products/:id
fn deleteProduct(r: zap.Request) void {
    if (r.getParamStr("id")) |id_str| {
        const id = std.fmt.parseInt(u32, id_str, 10) catch {
            r.setStatus(.bad_request) catch return;
            r.sendBody("Invalid ID") catch return;
            return;
        };

        for (products.items, 0..) |p, i| {
            if (p.id == id) {
                _ = products.orderedRemove(i);
                r.setStatus(.no_content) catch return;
                r.sendBody("") catch return;
                return;
            }
        }

        r.setStatus(.not_found) catch return;
        r.sendBody("Product not found") catch return;
    }
}

pub fn main() !void {
    // Add some sample products
    try products.append(.{
        .id = next_id,
        .name = "Laptop",
        .price = 999.99,
        .description = "High-performance laptop",
        .stock = 5,
    });
    next_id += 1;

    try products.append(.{
        .id = next_id,
        .name = "Mouse",
        .price = 29.99,
        .description = "Wireless mouse",
        .stock = 20,
    });
    next_id += 1;

    // Setup routes
    var listener = zap.HttpListener.init(.{
        .port = 3000,
        .on_request = null,
        .log = true,
    });
    defer listener.deinit();

    var router = zap.Router.init(gpa.allocator(), .{});
    defer router.deinit();

    // Register routes
    try router.handle_func("/products", getProducts, .GET);
    try router.handle_func("/products", createProduct, .POST);
    try router.handle_func("/products/:id", getProduct, .GET);
    try router.handle_func("/products/:id", deleteProduct, .DELETE);

    listener.on_request = router.on_request_handler();

    std.debug.print("Server running on http://localhost:3000\n", .{});
    std.debug.print("Endpoints:\n", .{});
    std.debug.print("  GET    /products     - List all products\n", .{});
    std.debug.print("  GET    /products/:id - Get product by ID\n", .{});
    std.debug.print("  POST   /products     - Create new product\n", .{});
    std.debug.print("  DELETE /products/:id - Delete product\n", .{});

    try listener.listen();

    zap.start(.{
        .threads = 2,
        .workers = 2,
    });
}