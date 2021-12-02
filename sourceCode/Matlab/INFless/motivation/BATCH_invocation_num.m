clear;
clc
load workload3_sub_max_70; 

util=zeros(1,length(workload))';
invocation_num=zeros(1,length(workload))';
timeout=200-40;
for i=1:length(workload);
    if workload(i)>0;
        interval=1000/workload(i);
        batch_length=ceil(timeout/interval);
        if batch_length<4;
            invocation_num(i)=ceil(workload(i)/batch_length);
            util(i)=batch_length/4;
        else
            invocation_num(i)=ceil(workload(i)/4);
            util(i)=1;
        end 
    else
        invocation_num(i)=0;
    end
end
plot(invocation_num)
hold on
plot(workload)