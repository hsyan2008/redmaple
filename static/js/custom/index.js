jQuery(document).ready(function(){

    ///// TRANSFORM CHECKBOX /////							
    jQuery('input:checkbox').uniform();

    ///// LOGIN FORM SUBMIT /////
    jQuery('#login').submit(function(){

        if(jQuery('#username').val() == '' || jQuery('#password').val() == '') {
            jQuery('.nousername').fadeIn();
            jQuery('.nousername .loginmsg').text("用户名或密码不能为空");
            return false;	
        }

        $.ajax({
            type: "POST",
            dataType: "json",
            url: "/index/login",
            data: $('#login').serialize(),
            success: function(msg){
                if (msg.err_no > 0)  {
                    jQuery('.nousername').fadeIn();
                    jQuery('.nousername .loginmsg').text(msg.err_msg);
                } else {
                    jQuery('.nousername').hide();
                    window.location = "/"
                }
            }
        });
        return false;
    });

    ///// ADD PLACEHOLDER /////
    jQuery('#username').attr('placeholder','Username');
    jQuery('#password').attr('placeholder','Password');
});
