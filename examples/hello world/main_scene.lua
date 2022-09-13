local ball

function on_load()
    -- sprite name (e.g., "Ball") must be unique in current scene
    ball = L.new_sprite("Ball")
    -- or:
    -- ball = Lemon.new_sprite("My Bunny")
    ball:set_texture_from_file("ball.png")
    ball.width, ball.height = 64, 64
    ball.x, ball.y = 400, 250
    ball:set_origin_center()
end

function on_update(dt)
    ball.x = ball.x + 400 * dt
    ball.rotation = ball.rotation + 500 * dt
    if ball.x > L.get_screen_width() then
        ball.x = 0
    end
end
