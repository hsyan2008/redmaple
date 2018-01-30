/*
 * 	Additional function for forms.html
 *	Written by ThemePixels	
 *	http://themepixels.com/
 *
 *	Copyright (c) 2012 ThemePixels (http://themepixels.com)
 *	
 *	Built for Amanda Premium Responsive Admin Template
 *  http://themeforest.net/category/site-templates/admin-templates
 */

jQuery(document).ready(function(){
	
	///// FORM TRANSFORMATION /////
	jQuery('input:checkbox, input:radio, select.uniformselect, input:file').uniform();


	///// FORM VALIDATION /////
	jQuery("#form1").validate({
		rules: {
			SmtpAddr: "required",
			SmtpPort: {
                required:true,
                range:[5,65535]
            },
			SmtpUser: "required"
		},
		messages: {
			SmtpAddr: "请输入smtp服务器地址",
			SmtpPort: {
                required:"请输入smtp端口",
                range:"端口必须在5-65535之间"
            },
			SmtpUser: "请输入smtp帐号"
		},
        submitHandler:function(form){
            $.ajax({
                type: "POST",
                dataType: "json",
                url: "/setting/save",
                data: $(form).serialize(),
                success: function(msg){
                    id = $("#Id").val()
                    if (id > 0) {
                        txt = "修改"
                    } else {
                        txt = "创建"
                    }
                    if (msg.err_no > 0)  {
                        jAlert(msg.err_msg, txt+'失败');
                    } else {
                        jAlert(txt+'成功', txt+'成功');
                    }
                }
            });
            return false;
        }
	});
	
	
	///// SPINNER /////
	
	jQuery("#spinner").spinner({min: 0, max: 100, increment: 2});
	
	///// SELECT WITH SEARCH /////
	jQuery(".chzn-select").chosen();
	
});
