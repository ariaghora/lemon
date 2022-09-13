local tex

function on_load()
    local scale = 1
    tex = L.new_texture("assets/basic_char_spritesheet.png")
    this:set_texture(tex)
    this.width = this.width * scale
    this.height = this.height * scale
    this.frame_width = 48
    this.frame_height = 48
    this.frame_count_x = 2
    this.frame_count_y = 1
    this.x = 200
    this.y = 100
    this:play()
    this:set_origin_center()
end

function on_update()
end
