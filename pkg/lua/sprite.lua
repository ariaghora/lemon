function __proto_new_sprite(sprite_name)
    local sprite = {}
    sprite.animated = false
    sprite.frame_index = 1
    sprite.frame_width = 0
    sprite.frame_height = 0
    sprite.frame_count_x = 1
    sprite.frame_count_y = 1
    sprite.animation_duration = 1
    sprite.frame_offset_x = 0
    sprite.frame_offset_y = 0
    sprite.name = sprite_name
    sprite.x = 0
    sprite.y = 0
    sprite.width = 0
    sprite.height = 0
    sprite.scale = 1
    sprite.rotation = 0
    sprite.origin_x = 0
    sprite.origin_y = 0
    sprite.playing = false
    sprite.shown = true
    sprite.play = function(sprite)
        sprite.playing = true
    end
    sprite.stop = function(sprite)
        sprite.playing = false
    end
    sprite.set_origin_center = function(self)
        self.origin_x = self.width / 2
        self.origin_y = self.height / 2
    end
    sprite.collides_with = function(other)
        return
    end
    -- sprite.set_texture(texture)
    -- sprite.set_texture_from_file(filename)
    -- sprite.remove()
    return sprite
end
