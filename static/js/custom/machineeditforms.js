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
			Name: "required",
            Env: {
                required: true,
                minlength: 1,
                min: 1,
            },
			Ip: "required",
			Port: "required",
			User: "required",
			Auth: "required",
			InnerIp: "required",
			InnerPort: "required",
			InnerUser: "required",
			InnerAuth: "required"
		},
		messages: {
			Name: "请输入服务器名称",
			Env: "请选择服务器环境",
			Ip: "请输入服务器IP地址",
			Port: "请输入服务器端口",
			User: "请输入服务器登陆名",
            Auth: "请输入认证信息",
			InnerIp: "请输入服务器IP地址",
			InnerPort: "请输入服务器端口",
			InnerUser: "请输入服务器登陆名",
            InnerAuth: "请输入认证信息"
		},
        submitHandler:function(form){
            $.ajax({
                type: "POST",
                dataType: "json",
                url: "/machine/save",
                data: $(form).serialize(),
                success: function(msg){
                    id = $("#Id").val()
                    if (id > 0) {
                        txt = "修改"
                    } else {
                        txt = "添加"
                    }
                    if (msg.err_no > 0)  {
                        jAlert(msg.err_msg, txt+'失败');
                    } else {
                        jAlert(txt+'成功', txt+'成功', function(){
                            window.location.href="/machine"
                        });
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
