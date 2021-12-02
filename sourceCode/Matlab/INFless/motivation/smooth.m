% clear
% clc 
clear new_data;
data=workload;
start=1;
interval=14;
row=1;
if interval==1;
    plot(data)
    return 
end

for i=1:length(data);
    if start+interval>length(data);
    break;
    end
    new_data(row,:)=mean(data(start:start+interval-1,:));
    start=start+interval;
    row=row+1;
end
stairs(new_data,'lineWidth',1.1)
 

%% ·Ö¶ÎÆ½»¬  
%  clc
% start=1;
% interval=30;
% row=1;
% if interval==1;
%     plot(data)
%     return 
% end
% 
% for i=1:length(data);
%     if start+interval>length(data);
%         break;
%     end
%     new_data(start:start+interval-1,1)=floor(mean(data(start:start+interval-1))+0.5);
%     start=start+interval;
%     row=row+1;
% end 
% hold on
% new_data=new_data(1:1100,1);
% plot(new_data(1:1100,1))
