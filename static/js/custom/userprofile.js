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


	///// DUAL BOX /////
	var db = jQuery('#dualselect').find('.ds_arrow .arrow');	//get arrows of dual select
	var sel1 = jQuery('#dualselect select:first-child');		//get first select element
	var sel2 = jQuery('#dualselect select:last-child');			//get second select element
	
	sel2.empty(); //empty it first from dom.
	
	db.click(function(){
		var t = (jQuery(this).hasClass('ds_prev'))? 0 : 1;	// 0 if arrow prev otherwise arrow next
		if(t) {
			sel1.find('option').each(function(){
				if(jQuery(this).is(':selected')) {
					jQuery(this).attr('selected',false);
					var op = sel2.find('option:first-child');
					sel2.append(jQuery(this));
				}
			});	
		} else {
			sel2.find('option').each(function(){
				if(jQuery(this).is(':selected')) {
					jQuery(this).attr('selected',false);
					sel1.append(jQuery(this));
				}
			});		
		}
	});
	
	
	
	///// FORM VALIDATION /////
	jQuery("#form1").validate({
		rules: {
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
            group_id: "required"
		},
		messages: {
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
                url: "/user/profile",
                data: $(form).serialize(),
                success: function(msg){
                    if (msg.err_no > 0)  {
                        jAlert(msg.err_msg, '修改失败');
                    } else {
                        jAlert('修改成功', '修改成功', function(){
                            window.location.href="/"
                        });
                    }
                }
            });
            return false;
        }
	});
	
	
	///// TAG INPUT /////
	
	jQuery('#tags').tagsInput();

	
	///// SPINNER /////
	
	jQuery("#spinner").spinner({min: 0, max: 100, increment: 2});
	
	
	///// CHARACTER COUNTER /////
	
	jQuery("#textarea2").charCount({
		allowed: 120,		
		warning: 20,
		counterText: 'Characters left: '	
	});
	
	
	///// SELECT WITH SEARCH /////
	jQuery(".chzn-select").chosen();
	
});
