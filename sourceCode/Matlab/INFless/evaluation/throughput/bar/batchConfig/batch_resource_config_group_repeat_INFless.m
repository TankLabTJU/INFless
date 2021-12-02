load resnet50_workload3_350msQPS;
data=resnet50_workload3_350msQPS;
repeat=zeros(1,20);
repeatCount=zeros(1,20);
count=0;
 
collect=zeros(20,5);
for i=1:length(data);
    score=data(i,1)+data(i,2)*100+data(i,3);
    locate=find(repeat==score);
    if ~isempty(locate); 
        repeatCount(locate)=repeatCount(locate)+1;
        continue;
    else
        data(i,:)
        count=count+1;
        repeat(count)=score;
        collect(count,1:4)=data(i,:);
        repeatCount(count)=repeatCount(count)+1;
    end
end
%% count and throuhgput efficiency
collect(:,6)=repeatCount';
for i=1:length(collect);
    collect(i,5)=collect(i,4)/(sum(collect(i,2)*64+collect(i,3)*142));
end
%% normalized to 1
collect(:,8)=collect(:,5)/max(collect(:,5));
collect(:,7)=collect(:,6)/max(collect(:,6));
