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