select 
  date,
  round(sum(suffer_score) over (
    order by date asc
    rows between
    unbounded preceding
    and current row
  ), 0) / 10,
  suffer_score / 10,
  (case when dist is null then 0 else round(dist / 1000, 1) end) 
from (
  select
    date,
    avg_speed,
    ( case avg_hr when 255.0 then 0 else avg_hr end) as avg_hr,
    ( case avg_hr when 255.0 then 0 else hhr end) as hhr,
    -- How much heart rate you have in reserve
    ( case avg_hr when 255.0 then 0 else 100 - hhr end) as hhr_left,
    -- How much VO2 you've used (see below)
    ( case avg_hr when 255.0 then 0 else hhr * 1.12 - 12 end) as vo2_max,  
    dist,
    time_min,
    alt_change,
    ( 
      case avg_hr 
        -- 255 rating means there was activity but no heart rate gathered
        -- use our best guess at suffer score
        when 255.0 then round((dist + alt_change + time_min) / 100 * 0.634, 2)
        -- -1 means there was no activity that day
        when -1 then -20
        -- otherwise we have some data do trimp_avg
        else round( (time_min * avg_hr) / 100, 0) 
      end 
    ) as suffer_score
  from (
    select
        dts.dte as date,
        round(avg( t.speed ), 2) as avg_speed,
        case when t.hr is null then -1
          else round(avg( t.hr ), 0)
        end as avg_hr,
        -- How much of your heart rate you've used
        round(avg( ((t.hr - h.rest_hr) / (h.run_max_hr - h.rest_hr)) * 100 ), 0) as hhr,
        max( t.dist ) as dist,
        round(max(t.alt) - min(t.alt), 2) as alt_change,
        -- This does odd things around midnight
        -- (max(strftime('%s', t.time)) - min(strftime('%s', t.time))) / 60 as time_min,
        round(sum(distinct l.total_time) / 60, 2) as time_min
    from dates_2019 dts
    join athlete h 
    left outer join activity a on strftime('%Y-%m-%d', a.time) = dts.dte
    left outer join trackpoint t on a.uuid = t.activity_uuid
    left outer join lap l on a.uuid = l.activity_uuid
    where dts.dte > '2019-09-28' and dts.dte <= '2019-12-31'
    group by dts.dte
    order by dts.dte asc
  )
)
