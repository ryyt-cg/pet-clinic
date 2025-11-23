


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
