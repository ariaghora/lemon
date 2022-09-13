function on_load()
    print("press space to start the game")
end

function on_update()
    L.draw_rect_fill(0, 0, L.get_screen_width(), L.get_screen_height(), L.RGBA(0, 0, 0, 255))
    L.draw_text("Press space to start", 240, L.get_screen_height() / 2, 30, L.RGBA(255, 255, 255, 255))
    if L.is_key_pressed(L.KEY_SPACE) then
        L.set_scene("main_scene")
    end
end
