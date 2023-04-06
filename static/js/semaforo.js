function fa_a_active() {
    $("#fa_v").removeClass("luz-v-active");
    $("#fa_v").addClass("luz-v");
    $("#fa_a").addClass("parpadear_a");
    return;
}

function fa_r_active() {
    $("#fa_a").removeClass("parpadear_a");
    $("#fa_r").removeClass("luz-r");
    $("#fa_r").addClass("luz-r-active");
    $("#fb_a").addClass("parpadear_a");
    setTimeout(fb_a_active, 0);
    return;
}

function fb_a_active() {
    $("#Fase_1").removeClass("active");
    $("#Fase_3").removeClass("active");
    $("#Fase_4").removeClass("active");
    $("#Fase_2").addClass("active");

    $("#fb_r").removeClass("luz-r-active");
    $("#fb_r").addClass("luz-r");
    $("#fb_a").addClass("parpadear_a");
    setTimeout(fb_v_active, f2_la());

    return;
}

function fb_v_active() {
    $("#fb_a").removeClass("parpadear_a");
    $("#fb_v").removeClass("luz-v");
    $("#fb_v").addClass("luz-v-active");
    setTimeout(fc_a_active, f2_lv());
    setTimeout(fb_restart, f2_lv());
    return;

}

function fb_restart() {
    $("#fb_v").removeClass("luz-v-active");
    $("#fb_v").addClass("luz-v");
    $("#fb_r").addClass("luz-r-active");
    return;
}

function fc_a_active() {
    $("#Fase_1").removeClass("active");
    $("#Fase_2").removeClass("active");
    $("#Fase_4").removeClass("active");
    $("#Fase_3").addClass("active");

    $("#fc_r").removeClass("luz-r-active");
    $("#fc_r").addClass("luz-r");
    $("#fc_a").addClass("parpadear_a");
    setTimeout(fc_v_active, f3_la());
    return;
}

function fc_v_active() {
    $("#fc_a").removeClass("parpadear_a");
    $("#fc_v").removeClass("luz-v");
    $("#fc_v").addClass("luz-v-active");
    setTimeout(fd_a_active, f3_lv());
    setTimeout(fc_restart, f3_lv());
    return;

}

function fc_restart() {
    $("#fc_v").removeClass("luz-v-active");
    $("#fc_v").addClass("luz-v");
    $("#fc_r").addClass("luz-r-active");
    return;
}

function fd_a_active() {
    $("#Fase_1").removeClass("active");
    $("#Fase_2").removeClass("active");
    $("#Fase_3").removeClass("active");
    $("#Fase_4").addClass("active");

    $("#fd_r").removeClass("luz-r-active");
    $("#fd_r").addClass("luz-r");
    $("#fd_a").addClass("parpadear_a");
    setTimeout(fd_v_active, f4_la());
    return;
}

function fd_v_active() {
    $("#fd_a").removeClass("parpadear_a");
    $("#fd_v").removeClass("luz-v");
    $("#fd_v").addClass("luz-v-active");
    setTimeout(fa_a_active, f4_lv());
    setTimeout(fd_restart, f4_lv());
    return;

}

function fd_restart() {
    $("#fd_v").removeClass("luz-v-active");
    $("#fd_v").addClass("luz-v");
    $("#fd_r").addClass("luz-r-active");
    $("#fa_r").removeClass("luz-r-active");
    $("#fa_r").addClass("luz-r");
    $("#fa_v").removeClass("luz-v");
    $("#fa_v").addClass("luz-v-active");
    $("#Fase_2").removeClass("active");
    $("#Fase_3").removeClass("active");
    $("#Fase_4").removeClass("active");
    $("#Fase_1").addClass("active");
    return;
}

