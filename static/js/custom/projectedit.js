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
    var db1 = jQuery('#dualselect1').find('.ds_arrow .arrow');	//get arrows of dual select
    var sel1 = jQuery('#dualselect1 select:first-child');		//get first select element
    var sel2 = jQuery('#dualselect1 select:last-child');			//get second select element
    // sel1.empty(); //empty it first from dom.
    sel1.find('option').attr('selected', true)
    db1.click(function(){
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

    ///// DUAL BOX /////
    var db2 = jQuery('#dualselect2').find('.ds_arrow .arrow');	//get arrows of dual select
    var sel3 = jQuery('#dualselect2 select:first-child');		//get first select element
    var sel4 = jQuery('#dualselect2 select:last-child');			//get second select element
    // sel3.empty(); //empty it first from dom.
    sel3.find('option').attr('selected', true)
    db2.click(function(){
        var t = (jQuery(this).hasClass('ds_prev'))? 0 : 1;	// 0 if arrow prev otherwise arrow next
        if(t) {
            sel3.find('option').each(function(){
                if(jQuery(this).is(':selected')) {
                    jQuery(this).attr('selected',false);
                    var op = sel4.find('option:first-child');
                    sel4.append(jQuery(this));
                }
            });	
        } else {
            sel4.find('option').each(function(){
                if(jQuery(this).is(':selected')) {
                    jQuery(this).attr('selected',false);
                    sel3.append(jQuery(this));
                }
            });		
        }
        sel3.find('option').each(function(){
            jQuery(this).attr('selected',true);
        });
    });

    ///// DUAL BOX /////
    var db3 = jQuery('#dualselect3').find('.ds_arrow .arrow');	//get arrows of dual select
    var sel5 = jQuery('#dualselect3 select:first-child');		//get first select element
    var sel6 = jQuery('#dualselect3 select:last-child');			//get second select element
    // sel5.empty(); //empty it first from dom.
    sel5.find('option').attr('selected', true)
    db3.click(function(){
        var t = (jQuery(this).hasClass('ds_prev'))? 0 : 1;	// 0 if arrow prev otherwise arrow next
        if(t) {
            sel5.find('option').each(function(){
                if(jQuery(this).is(':selected')) {
                    jQuery(this).attr('selected',false);
                    var op = sel6.find('option:first-child');
                    sel6.append(jQuery(this));
                }
            });	
        } else {
            sel6.find('option').each(function(){
                if(jQuery(this).is(':selected')) {
                    jQuery(this).attr('selected',false);
                    sel5.append(jQuery(this));
                }
            });		
        }
        sel5.find('option').each(function(){
            jQuery(this).attr('selected',true);
        });
    });

    ///// FORM VALIDATION /////
    // Smart Wizard 	
    id = $("#Id").val()
    if (id > 0) {
        enableAllSteps = true
    } else {
        enableAllSteps = false
    }
    jQuery('#wizard').smartWizard({
        onFinish: onFinishCallback,
        labelNext:"下一步",
        labelPrevious:"上一步",
        labelFinish:"提交",
        enableAllSteps:enableAllSteps,    //编辑的时候
    });

    function onFinishCallback(){
        //这样会调用validate
        $("#form1").submit()
    } 

    jQuery(".inline").colorbox({inline:true, width: '60%', height: '500px'});

    jQuery('select, input:checkbox').uniform();

    ///// FORM VALIDATION /////
    jQuery("#form1").validate({
        ignore:"",
        rules: {
            Name: "required",
            Git: "required",
            Wwwroot: "required",
            //devMachineIds: "required",
            testMachineIds: "required",
            prodMachineIds: "required",
        },
        messages: {
            Name: "请输入项目名称",
            Git: "请输入git地址",
            Wwwroot: "请输入代码部署路径",
            //devMachineIds: "请选择开发环境服务器",
            testMachineIds: "请选择测试环境服务器",
            prodMachineIds: "请选择生产环境服务器",
        },
        errorContainer: ".notibar",
        errorLabelContainer: ".notibar p",
        wrapper: "div",
        submitHandler:function(form){
            $.ajax({
                type: "POST",
                dataType: "json",
                url: "/project/save",
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
                        jAlert(txt+'成功', txt+'成功', function(){
                            window.location.href="/project"
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
