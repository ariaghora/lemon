# Sprite

A sprite is represented as a lua table.
It is instantiated with `L.new_sprite()` method.

## Attributes

| Attribute  | Description                                                                                                                     |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `x`        | Horizontal position of a sprite                                                                                                 |
| `y`        | Vertical position of a sprite                                                                                                   |
| `rotation` | Rotation amount of a sprite. Positive increment is for clockwise rotation, negative increment is for counter-clockwise rotation |
| `scale`    | Sprite scale, defaults to 0                                                                                                     |

## Methods

- `set_texture(texture_data)`
- `set_texture_from_file(file_name)`