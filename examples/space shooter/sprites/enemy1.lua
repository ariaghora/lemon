local velocity = 300
local explosion_texs = {}
local current_explosion_tex = nil
local explosion_counter = 0
local explosion_idx = 1

function on_load()
    this.x = L.get_screen_width()
    this.y = math.random(this.height / 2, L.get_screen_height() - this.height / 2)

    for i = 1, 8 do
        explosion_texs[i] = L.new_texture("assets/images/explosion/" .. tostring(i) .. ".png")
    end
end

function on_update(dt)
    this.x = this.x - velocity * dt

    if this.x < -this.width then
        this.remove()
    end

    if this.destroyed then
        this.shown = false
        velocity = velocity - 10
        L.draw_texture(explosion_texs[explosion_idx], this.x - 128, this.y - 128)

        explosion_counter = explosion_counter + dt
        if explosion_counter > 0.05 then
            explosion_counter = 0

            explosion_idx = explosion_idx + 1
            if explosion_idx > 8 then
                this.remove()
            end
        end
    end
end
