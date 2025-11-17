const std = @import("std");


pub fn main() !void {
    // const allocator = gpa.allocator();
    // defer _ = gpa.deinit();
    //
    // // Load configuration
    // var config = try Config.load(allocator, "config.json");
    // defer config.deinit(allocator);
    //
    // std.debug.print("Configuration loaded:\n", .{});
    // std.debug.print("  Server: {s}:{d}\n", .{ config.server.host, config.server.port });
    // std.debug.print("  Database: {s}\n", .{config.database.path});
    //
    // // Initialize database
    // db = try Database.init(allocator, config.database.path);
    // defer db.deinit();
    // try db.createTables();
    //
    // // Initialize repository, handler and dependency injection
    // repo = ProductRepository.init(allocator, &db);
    // handler = ProductHandler.init(allocator, &repo);
    //
    // // Setup Zap server
    // var listener = zap.HttpListener.init(.{
    //     .port = config.server.port,
    //     .on_request = on_request,
    //     .log = true,
    // });
    // try listener.listen();

    // std.debug.print("\n🚀 Server running on http://{s}:{d}\n", .{ config.server.host, config.server.port });
    std.debug.print("\nAvailable endpoints:\n", .{});
    std.debug.print("  GET    /health              - Health check\n", .{});
    std.debug.print("  GET    /api/products        - Get all products\n", .{});
    std.debug.print("  GET    /api/products/:id    - Get product by ID\n", .{});
    std.debug.print("  POST   /api/products        - Create new product\n", .{});
    std.debug.print("  DELETE /api/products/:id    - Delete product\n", .{});

    // zap.start(.{
    //     .threads = 4,
    //     .workers = 2,
    // });
}