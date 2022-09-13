-- count how many seconds passed since the game scene started
-- we can control game flow (which enemy will appear, which background
-- will appear, etc.) according to game_elapsed_seconds
local game_elapsed_seconds = 0

local enemies_per_second = 1
local enemy_counter = 0
local enemy_texture = nil

local game_counter = 0

local num_stars = 150
local stars = {}

function on_load()
    L.set_global("score", 0)
    -- create our player and set its controller script
    local player = L.new_sprite("player")
    player:set_script("sprites/player.lua")

    -- create stars
    for i = 1, num_stars do
        stars[i] = {}
        stars[i].x = math.random(0, L.get_screen_width())
        stars[i].y = math.random(0, L.get_screen_height())
    end

    -- create enemy_texture from file
    enemy_texture = L.new_texture("assets/images/ship2.png")
end

function on_update(dt)
    -- clear our screen with black color
    L.draw_rect_fill(0, 0, L.get_screen_width(), L.get_screen_height(), L.RGBA(0, 0, 0, 255))

    for i = 1, #stars do
        L.draw_rect_fill(stars[i].x, stars[i].y, 3, 3, L.RGBA(255, 255, 255, 100))
        stars[i].x = stars[i].x - 20 * dt
        if stars[i].x < -2 then
            stars[i].x = L.get_screen_width()
        end
    end

    -- print out the score on screen
    L.draw_text(
        "Score: " .. tostring(L.get_global("score")),
        10,
        10,
        20,
        L.RGBA(255, 255, 255, 255)
    )

    -- Generate enemies randompy at a rate of enemies_per_second
    if game_counter % math.floor(L.get_fps() / enemies_per_second) == 0 then
        enemy_counter = enemy_counter + 1
        local e = L.new_sprite("enemy" .. tostring(enemy_counter))
        -- e:set_texture_from_file("assets/images/ship2.png")
        e:set_texture(enemy_texture)
        e:set_script("sprites/enemy1.lua")
    end

    -- increase elapsed seconds
    if game_counter > L.get_fps() then
        game_elapsed_seconds = game_elapsed_seconds + 1
        game_counter = 0
    end

    -- increment game counter
    game_counter = game_counter + 1
end
