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
	
	// sel1.empty(); //empty it first from dom.
	
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
        sel1.find('option').each(function(){
            jQuery(this).attr('selected',true);
        });
	});

	///// FORM VALIDATION /////
	jQuery("#form1").validate({
		rules: {
			Name: "required",
			Comment: "required",
			ReviewUserId: "required",
            projectIds: {
                required: true,
            }
		},
		messages: {
			Name: "请输入本任务名称",
			Comment: "请输入开发内容",
			ReviewUserId: "请选择review人员",
			projectIds: "请选择开发分支"
		},
        submitHandler:function(form){
            $.ajax({
                type: "POST",
                dataType: "json",
                url: "/task/save",
                data: $(form).serialize(),
                success: function(msg){
                    if (msg.err_no > 0)  {
                        jAlert(msg.err_msg, '创建失败');
                    } else {
                        jAlert('创建成功', '创建成功', function(){
                            window.location.href="/task"
                        });
                    }
                }
            });
            return false;
        }
	});
	
	
	///// SELECT WITH SEARCH /////
	jQuery(".chzn-select").chosen();
	
});
