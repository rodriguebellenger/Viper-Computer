#include <chrono>
#include <cmath>
#include <print>
#include <cstdint>
#include <cstdlib>

#include "opcode.h"

using s8  = std::int8_t;
using s16 = std::int16_t;
using s32 = std::int32_t;
using s64 = std::int64_t;

using u8  = std::uint8_t;
using u16 = std::uint16_t;
using u32 = std::uint32_t;
using u64 = std::uint64_t;

using f32 = float;
using f64 = double;

u8 RAM[1024] = {0};

int main() {
    RAM[0] = 200;
    std::println("{}", RAM[0]);
    std::println("{}", RAM[1]);

    for (const auto& [opcode, types] : opNames) {
        std::print("{} : ", opcode);
        for (const auto& type : types) {
            std::string str = "";
            switch (type) {
                case Null:
                    str = "";
                    break;
                case Register:
                    str = "Register";
                    break;
                case Address:
                    str = "Address";
                    break;
                case Size:
                    str = "Size";
                    break;
                case Offset:
                    str = "Offset";
                    break;
                case Comparison:
                    str = "Comparison";
                    break;
                case Int8:
                    str = "Int8";
                    break;
                case Int16:
                    str = "Int16";
                    break;
                case Int32:
                    str = "Int32";
                    break;
            }

            std::print("{} ", str);
        }
        std::println();
    }

    return 0;
}
