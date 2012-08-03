/** 微博时间格式化显示
 * @param $timestamp，标准时间戳
 */
function smc_time_since($timestamp) {
    $since = abs(time()-$timestamp);
    $gmt_offset = get_option('gmt_offset') * 3600;//获取wordpress的时区偏移值
    $timestamp += $gmt_offset; $current_time = mktime() + $gmt_offset;
    if(floor($since/3600)){
        if(gmdate('Y-m-d',$timestamp) == gmdate('Y-m-d',$current_time)){
            $output = '今天 ';
            $output.= gmdate('H:i',$timestamp);
        }else{
            if(gmdate('Y',$timestamp) == gmdate('Y',$current_time)){
                $output = gmdate('m月d日 H:i',$timestamp);
            }else{
                $output = gmdate('Y年m月d日 H:i',$timestamp);
            }
        }
    }else{
        if(($output=floor($since/60))){
            $output = $output.'分钟前';
        }else $output = '刚刚';
    }
    return $output;
}