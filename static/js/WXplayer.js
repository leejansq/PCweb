

       function StarTimer(o){
		getInterval(o);
        //setInterval("getInterval(o)",100000);//1000为1秒钟
       }
       var list_url=[]; 
	   var story_url;
	var play_flag=false;
	
	function state_juge(os,index){
		setTimeout(function(){
			//alert(os.networkState);
			if(os.networkState==3){
				os.src=list_url[index];
				os.play();
				state_juge(os,index)
			}
		},1000);
	}
       function getInterval(ob)
       {
        var   url='/v2/api/m3u8'
        var jsonarr = {'uid':eid};
				$.ajax({
				type:'get',
				url:url,
				data:jsonarr,
				timeout:90000,
				cache:false,
//				beforeSend:function(){
//						Win.Loading();			
//				},
				dataType : 'json',
				success:function(o){
				
						//Win.Loading("CLOSE"); 
					if(o.Code == "20000"){
					    //alert(o.Data[0].Url);
						var html='';
						list_url=[];
						 $.each(o.Data, function(commentIndex, comment){
                               html +=comment.Url+"\n";
							if(play_flag==true){
								if(comment.Url>story_url){
									list_url.push(comment.Url);
								}
							}else{
								list_url.push(comment.Url);
							}
							
                         });
						alert(list_url);
						alert(html);
						if(list_url.length>0){
							story_url=list_url[list_url.length-1];
						}
						var index=0;
						if(play_flag==true) return;
						play_flag=true;
						//alert(list_url[index])
						ob.src=list_url[index++];
						
						ob.play();
						ob.addEventListener("ended",function(){
							if(list_url.length==0) {
								alert("No New Video Attached!")
								return;
							}
							if(index==list_url.length-1) {
								getInterval(ob);
								//alert(index);
								ob.src=list_url[index++];
								ob.play();
								state_juge(ob,index-1);
								index=0;
							}else{
								//alert(index);
								ob.src=list_url[index++];
								ob.play();
								state_juge(ob,index-1);
								//alert(ob.networkState)
							}
							
						})
						
					}else{
						alert("get_url error!");	
					}
				}
			})
            //alert('aaaaaaaaa');
        }
		