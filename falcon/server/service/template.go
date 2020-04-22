package service

func GetProMonitorTemplate() string {
	return `
		<![CDATA[
        <tr class='<%=(alternate ? "dark" : "")%> <%=(this.pid == "" ? "total" : "")%>'>
            <td class="name" title='<%= (this.pid == "" ? "" : this.name) %>'><%= (this.pid == "" ? "<br/>" : this.name) %></td>
            <td class="name" title='<%= (this.pid == "" ? "" : this.bit) %>'><%= (this.pid == "" ? "" : this.bit) %></td>
            <td class="name" title='<%= (this.pid == "" ? "" : this.user_name) %>'><%= (this.pid == "" ? "" : this.user_name) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : this.pid) %>'><%= (this.pid == "" ? "" : this.pid) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : this.ppid) %>'><%= (this.pid == "" ? "" : this.ppid) %></td>
            <td class="name" title='<%= (this.pid == "" ? "" : formatSeconds(this.run_time) + " (" + this.run_time + "s)" ) %>'><%= (this.pid == "" ? "" : formatSeconds(this.run_time) + " (" + this.run_time + "s)" ) %></td>
            <td class="name" title='<%= (this.pid == "" ? "" : this.cmd.replace(/\'/g, "$#39;")) %>'><%= (this.pid == "" ? "" : this.cmd) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : this.handles) %>'><%= (this.pid == "" ? "" : this.handles) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : this.threads) %>'><%= (this.pid == "" ? "" : this.threads) %></td>
            <td class="name" title='<%= (this.pid == "" ? "" : formatSeconds(this.cpu_time) + " (" + this.cpu_time + "s)" ) %>'><%= (this.pid == "" ? "" : formatSeconds(this.cpu_time) + " (" + this.cpu_time + "s)" ) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : this.cpu_percent) %>'><%= (this.pid == "" ? "" : this.cpu_percent) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : ((this.rss)/1024/1024).toFixed(2)) %>'><%= (this.pid == "" ? "" : ((this.rss)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : ((this.vms)/1024/1024).toFixed(2)) %>'><%= (this.pid == "" ? "" : ((this.vms)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : ((this.io_read_bytes_diff)/1024/1024).toFixed(2)) %>'><%= (this.pid == "" ? "" : ((this.io_read_bytes_diff)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : ((this.io_write_bytes_diff)/1024/1024).toFixed(2)) %>'><%= (this.pid == "" ? "" : ((this.io_write_bytes_diff)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : this.io_read_count_diff) %>'><%= (this.pid == "" ? "" : this.io_read_count_diff) %></td>
            <td class="numeric" title='<%= (this.pid == "" ? "" : this.io_write_count_diff) %>'><%= (this.pid == "" ? "" : this.io_write_count_diff) %></td>
        </tr>
        <% alternate = !alternate; %>
        ]]>
	`
}

func GetSysMonitorTemplate() string {
	return `
        <![CDATA[
        <tr class='<%=(alternate ? "dark" : "")%> <%=(this.last_time == "" ? "total" : "")%>'>
            <td class="name" title='<%= (this.last_time == "" ? "" : this.last_time) %>'><%= (this.last_time == "" ? "<br/>" : this.last_time) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : this.processes) %>'><%= (this.last_time == "" ? "" : this.processes) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : this.handles) %>'><%= (this.last_time == "" ? "" : this.handles) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : this.threads) %>'><%= (this.last_time == "" ? "" : this.threads) %></td>
			<td class="name" title='<%= (this.last_time == "" ? "" : formatSeconds(this.cpu_time) + " (" + this.cpu_time + "s)" ) %>'><%= (this.last_time == "" ? "" : formatSeconds(this.cpu_time) + " (" + this.cpu_time + "s)" ) %></td>
			<td class="numeric" title='<%= (this.last_time == "" ? "" : this.cpu_percent) %>'><%= (this.last_time == "" ? "" : this.cpu_percent) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : ((this.rss)/1024/1024).toFixed(2)) %>'><%= (this.last_time == "" ? "" : ((this.rss)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : ((this.vms)/1024/1024).toFixed(2)) %>'><%= (this.last_time == "" ? "" : ((this.vms)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : ((this.net_bytes_recv_diff)/1024/1024).toFixed(2)) %>'><%= (this.last_time == "" ? "" : ((this.net_bytes_recv_diff)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : ((this.net_bytes_sent_diff)/1024/1024).toFixed(2)) %>'><%= (this.last_time == "" ? "" : ((this.net_bytes_sent_diff)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : this.net_packets_recv_diff) %>'><%= (this.last_time == "" ? "" : this.net_packets_recv_diff) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : this.net_packets_sent_diff) %>'><%= (this.last_time == "" ? "" : this.net_packets_sent_diff) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : ((this.disk_used)/1024/1024/1024).toFixed(2)) %>'><%= (this.last_time == "" ? "" : ((this.disk_used)/1024/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : ((this.io_read_bytes_diff)/1024/1024).toFixed(2)) %>'><%= (this.last_time == "" ? "" : ((this.io_read_bytes_diff)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : ((this.io_write_bytes_diff)/1024/1024).toFixed(2)) %>'><%= (this.last_time == "" ? "" : ((this.io_write_bytes_diff)/1024/1024).toFixed(2)) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : this.io_read_count_diff) %>'><%= (this.last_time == "" ? "" : this.io_read_count_diff) %></td>
            <td class="numeric" title='<%= (this.last_time == "" ? "" : this.io_write_count_diff) %>'><%= (this.last_time == "" ? "" : this.io_write_count_diff) %></td>
        </tr>
        <% alternate = !alternate; %>
        ]]>
	`
}

func GetOsInfoTemplate() string {
	return `
        <![CDATA[
        <tr class='<%=(alternate ? "dark" : "")%> <%=(this.host_id == "" ? "total" : "")%>'>
            <td class="name" title='<%= (this.host_id == "" ? "" : "Start: " + this.start_time + "&#10;Host: " + this.host_name + "&#10;Os: " + this.os_type + " (" + this.bit + " bit)&#10;User: " + this.process_user_name + "&#10;Pid: " + this.process_pid) %>'><%= (this.host_id == "" ? "<br/>" : "Start: " + this.start_time + "<br/>Host: " + this.host_name + "<br/>Os: " + this.os_type + " (" + this.bit + " bit)<br/>User: " + this.process_user_name + "<br/>Pid: " + this.process_pid) %></td>
            <td class="name" title='<%= (this.host_id == "" ? "" : "Name: " + this.cpu_model_name + "&#10;Speed: " + this.cpu_mhz + " MHz&#10;Counts: " + this.cpu_count + "&#10;Threads: " + this.cpu_threads) %>'><%= (this.host_id == "" ? "" : "Name: " + this.cpu_model_name + "<br/>Speed: " + this.cpu_mhz + " MHz<br/>Counts: " + this.cpu_count + "<br/>Threads: " + this.cpu_threads) %></td>
            <td class="name" title='<%= (this.host_id == "" ? "" : "Virtual: " + Math.ceil((this.men_virtual_total)/1024/1024/1024) + " GB&#10;Swap: " + Math.ceil((this.mem_swap_total)/1024/1024/1024) + " GB") %>'><%= (this.host_id == "" ? "" : "Virtual: " + Math.ceil((this.men_virtual_total)/1024/1024/1024) + " GB<br/>Swap: " + Math.ceil((this.mem_swap_total)/1024/1024/1024) + " GB") %></td>
            <td class="name" title='<%= (this.host_id == "" ? "" : Math.ceil((this.disk_size_total)/1024/1024/1024) + " GB") %>'><%= (this.host_id == "" ? "" : Math.ceil((this.disk_size_total)/1024/1024/1024) + " GB") %></td>
            <td class="name" title='<%= (this.host_id == "" ? "" : "Name : " + this.net_card_name + "&#10;Mac: " + this.mac_address + "&#10;Mask: " + this.mask_address + "&#10;Local IP: " + this.ip_address + "&#10;Public IP: " + this.public_address  + "&#10;Speed: " + this.net_speed + " Mb") %>'><%= (this.host_id == "" ? "" : "Name : " + this.net_card_name + "<br/>Mac: " + this.mac_address + "<br/>Mask: " + this.mask_address + "<br/>Local IP: " + this.ip_address + "<br/>Public IP: " + this.public_address  + "<br/>Speed: " + this.net_speed + " Mb") %></td>
        </tr>
        <% alternate = !alternate; %>
        ]]>
	`
}