const std = @import("std");

pub const TokenType = enum {
    // Keywords
    Keyword_BeginProject,
    Keyword_EndProject,
    Keyword_Setup,
    Keyword_Run,
    Keyword_Define,
    Keyword_Const,
    Keyword_Integer,
    Keyword_Float,
    // Operators
    OpAssign,
    OpPlus,
    OpMinus,
    OpMult,
    OpDiv,
    Semicolon,
    Colon,
    Pipe,
    LeftBrace,
    RightBrace,
    // Literals
    Identifier, 
    IntLiteral,
    FloatLiteral,
    EOF,
};

pub const Token = struct {
    tag: TokenType,
    slice: []const u8, // The actual val
    line: usize,
    column: usize,
};