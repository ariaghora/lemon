local velocity = 700
local bullet_tex

function on_load()
    bullet_tex = L.new_texture("assets/images/player_bullet.png")
    this:set_texture(bullet_tex)
end

function on_update(dt)
    this.x = this.x + velocity * dt

    -- Recall that all our enemies are created with name "enemy1", "enemy2", etc.
    -- We can find all sprites with name containing "enemy" and list them as a
    -- table which we can iterate over.
    local enemies = L.find_sprites_by_name_like("enemy")

    for i = 1, #enemies, 1 do
        local enemy = enemies[i]
        -- Detect simple collision when bullet center enters enemy's
        -- rect area. Remove both enemy and this bullet from current scene
        -- when collosion is detected.
        if ((this.x >= enemy.x - enemy.origin_x) and
            (this.y >= enemy.y - enemy.origin_y) and
            (this.x <= enemy.x - enemy.origin_x + enemy.width) and
            (this.y <= enemy.y - enemy.origin_y + enemy.height)) then
            if not enemy.destroyed then
                enemy.destroyed = true
                this.remove()
                L.set_global("score", L.get_global("score") + 1)
            end
        end
    end

    if this.x > L.get_screen_width() then
        this.remove()
    end
end
