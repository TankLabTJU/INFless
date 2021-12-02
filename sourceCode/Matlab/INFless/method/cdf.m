function[] =cdf(x,num)
xi = linspace(0,max(x),num);
F = ksdensity(x,xi,'function','cdf');
plot(xi,F);
xlim([min(x),max(x)]);
end