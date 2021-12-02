 clear
 clc
% clear instance_scale1;
% load instance_scale1;  
% data=instance_scale1(:,1);
clear data;
load requestNum;
data=requestNum(1:3006,3);
start=1;
interval=6;
row=1;
if interval==1;
    plot(data)
    return 
end

for i=1:length(data);
    if start+interval>length(data);
    break;
    end
    new_data(row,:)=floor(mean(data(start:start+interval-1,:)));
    start=start+interval;
    row=row+1;
end
plot(new_data,'lineWidth',1.1)
 