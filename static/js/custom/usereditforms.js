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
			name: "required",
			realname: "required",
			email: {
				required: true,
				email: true,
			},
			password: {
				minlength: 5,
			},
			repassword: {
                equalTo: "#password"
            },
            group_id: {
                required: true,
                minlength: 1,
                min: 1,
            }
		},
		messages: {
			name: "请输入帐号",
			realname: "请输入真实姓名",
			email: "请输入正确的email地址",
			password: {
                minlength: "密码长度不能小于 5 个字母"
            },
			repassword: {
                equalTo: "两次密码输入不一致"
            },
			group_id: "请选择分组"
		},
        submitHandler:function(form){
            $.ajax({
                type: "POST",
                dataType: "json",
                url: "/user/save",
                data: $(form).serialize(),
                success: function(msg){
                    id = $("#id").val()
                    if (id > 0) {
                        txt = "修改"
                    } else {
                        txt = "创建"
                    }
                    if (msg.err_no > 0)  {
                        jAlert(msg.err_msg, txt+'失败');
                    } else {
                        jAlert(txt+'成功', txt+'成功', function(){
                            window.location.href="/user"
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
