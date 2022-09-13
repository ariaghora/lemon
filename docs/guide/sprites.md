# Sprite

## Shared texture
It is not efficient for sprites with same texture to load texture data from file separately.
It will only fill the memory with redundant data.
Better approach: we can preload the texture data beforehand and let different sprites access it in the same time.

## Sprite script
Modular codes are always better.

## Animated sprite (via spritesheet)