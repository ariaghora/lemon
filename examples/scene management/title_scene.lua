function on_update(dt)
    L.draw_rect_fill(
        0,
        0,
        L.get_screen_width(),
        L.get_screen_height(),
        L.rgba(0, 0, 0, 255)
    )

    L.draw_text(
        "This is title scene\n" ..
        "Press spacebar to go to the main scene",
        0,
        0,
        30,
        L.rgba(255, 255, 255, 255)
    )

    if L.is_key_pressed(L.KEY_SPACE) then
        L.set_scene("main_scene")
    end
end
