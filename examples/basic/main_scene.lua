local droid

local droid_run_texture
local droid_death_texture
local droid_dead = false

local rabbito

function on_load()
    droid = L.new_sprite("droid")
    droid_run_texture = L.new_texture("assets/run.png")
    droid_death_texture = L.new_texture("assets/damaged and death.png")
    droid:set_texture(droid_run_texture)
    droid.x, droid.y, droid.width, droid.height = 200, 300, 58 * 5, 41 * 5
    droid.frame_count_x = 1
    droid.frame_count_y = 6
    droid.frame_height = 41
    droid.animation_duration = 0.3 --seconds to complete one animation_duration loop

    -- you can "wrap" sprite-related script in a separated lua file.
    -- Check rabbito.lua
    rabbito = L.new_sprite("rabbito")
    rabbito:set_script("rabbito.lua")
end

function on_update(dt)
    L.draw_rect_fill(
        0,
        0,
        L.get_screen_width(),
        L.get_screen_height(),
        L.RGBA(240, 221, 192, 255)
    )
    L.draw_rect_fill(
        600,
        0,
        20,
        L.get_screen_height(),
        L.RGBA(255, 255, 255, 255)
    )
    if not droid_dead then
        if L.is_key_down(L.KEY_RIGHT) then
            droid:play()
            droid.x = droid.x + 400 * dt
        elseif L.is_key_down(L.KEY_LEFT) then
            droid:play()
            droid.x = droid.x - 400 * dt
        else
            droid.frame_index = 1
            droid:stop()
        end
    end

    if droid.x > 600 then
        droid_dead = true
    end

    if droid_dead then
        droid:set_texture(droid_death_texture)
        droid.animation_duration = 0.8
        droid.frame_count_x = 1
        droid.frame_count_y = 8
        droid.width, droid.height = 58 * 5, 41 * 5
        droid.frame_height = 41
        if droid.frame_index == 8 then
            droid:stop()
            droid.frame_index = 8
        end
    end
end
