package udger

const (
	CrawlerClassId    = 99
	CrawlerClass     = "Crawler"
	CrawlerCode      = "crawler"
	BotDetailUrl     = "https://udger.com/resources/ua-list/bot-detail?bot="
	OsDetailUrl      = "https://udger.com/resources/ua-list/os-detail?os="
	BrowserDetailUrl = "https://udger.com/resources/ua-list/browser-detail?browser="
	DeviceDetailUrl  = "https://udger.com/resources/ua-list/device-detail?device="
	IpV4 = "v4"
	IpV6 = "v6"

	BrowsersSqlQuery = `SELECT
			  udger_client_class.id as class_id,
			  udger_client_regex.client_id as client_id,
			  regstring,
			  name,
			  name_code,
			  homepage,
			  icon,
			  icon_big,
			  engine,
			  vendor,
			  vendor_code,
			  vendor_homepage,
			  uptodate_current_version,
			  client_classification,
			  client_classification_code
			FROM udger_client_regex
			  JOIN udger_client_list ON udger_client_list.id = udger_client_regex.client_id
			  JOIN udger_client_class ON udger_client_class.id = udger_client_list.class_id
			ORDER BY sequence
			  ASC`

	OsSqlQuery = `SELECT
			  os_id,
			  regstring,
			  family,
			  family_code,
			  name,
			  name_code,
			  homepage,
			  icon,
			  icon_big,
			  vendor,
			  vendor_code,
			  vendor_homepage
			FROM udger_os_regex
			  JOIN udger_os_list ON udger_os_list.id = udger_os_regex.os_id
			ORDER BY sequence
			  ASC`

	OsRelationSqlQuery = `SELECT
			  os_id,
			  family,
			  family_code,
			  name,
			  name_code,
			  homepage,
			  icon,
			  icon_big,
			  vendor,
			  vendor_code,
			  vendor_homepage
			FROM udger_client_os_relation
			  JOIN udger_os_list ON udger_os_list.id = udger_client_os_relation.os_id
			WHERE client_id = %v`

	DevicesSqlQuery = `SELECT
			  deviceclass_id,
			  regstring,
			  name,
			  name_code,
			  icon,
			  icon_big
			FROM udger_deviceclass_regex
			  JOIN udger_deviceclass_list ON udger_deviceclass_list.id = udger_deviceclass_regex.deviceclass_id
			ORDER BY sequence
			  ASC`

	DeviceClassListSqlQuery = `SELECT
		udger_client_class.id,
		deviceclass_id,
		name,
		name_code,
		icon,
		icon_big
	   FROM udger_deviceclass_list
	   JOIN udger_client_class ON udger_client_class.deviceclass_id=udger_deviceclass_list.id`

	DeviceCodeListSqlQuery = `SELECT
	   id,
       os_family_code,
	   os_code,
	   regstring
	  FROM
	   udger_devicename_regex`

	BrandByRegexIdAndCodeSqlQuery = `SELECT
	  marketname,
	  brand,
	  brand_code,
	  brand_url,
	  icon,
	  icon_big
	FROM udger_devicename_list
	  JOIN udger_devicename_brand ON udger_devicename_brand.id = udger_devicename_list.brand_id
	WHERE regex_id=%v and code='%s'`

	IpSqlQuery = `SELECT
				 udger_ip_list.ip as ip,
   				 udger_crawler_list.id as botid,
   				 ip_last_seen,
   				 ip_hostname,
   				 ip_country,
   				 ip_city,
   				 ip_country_code,
   				 ip_classification,
   				 ip_classification_code,
   				 name,
   				 ver,
   				 ver_major,
   				 last_seen,
   				 respect_robotstxt,
   				 family,
   				 family_code,
   				 family_homepage,
   				 family_icon,
   				 vendor,
   				 vendor_code,
   				 vendor_homepage,
   				 crawler_classification,
   				 crawler_classification_code
			FROM udger_ip_list
				JOIN udger_ip_class ON udger_ip_class.id = udger_ip_list.class_id
				LEFT JOIN udger_crawler_list ON udger_crawler_list.id = udger_ip_list.crawler_id
				LEFT JOIN udger_crawler_class ON udger_crawler_class.id = udger_crawler_list.class_id
			ORDER BY sequence`

	CrawlersSqlQuery = `SELECT
		  udger_crawler_list.id as id,
		  name,
		  ver,
		  ver_major,
		  last_seen,
		  respect_robotstxt,
		  family,
		  family_code,
		  family_homepage,
		  family_icon,
		  vendor,
		  vendor_code,
		  vendor_homepage,
		  crawler_classification,
		  crawler_classification_code,
		  ua_string
		FROM udger_crawler_list
		  LEFT JOIN udger_crawler_class ON udger_crawler_class.id=udger_crawler_list.class_id`

	CrawlersByUaSqlQuery = `SELECT
		  udger_crawler_list.id as id,
		  name,
		  ver,
		  ver_major,
		  last_seen,
		  respect_robotstxt,
		  family,
		  family_code,
		  family_homepage,
		  family_icon,
		  vendor,
		  vendor_code,
		  vendor_homepage,
		  crawler_classification,
		  crawler_classification_code,
		  ua_string
		FROM udger_crawler_list
		  LEFT JOIN udger_crawler_class ON udger_crawler_class.id=udger_crawler_list.class_id
        where ua_string = '%s'`

	GetDataCenterV4SqlQuery = `select
		  name,
		  name_code,
		  homepage
		FROM udger_datacenter_range
		JOIN udger_datacenter_list ON udger_datacenter_range.datacenter_id=udger_datacenter_list.id
		  where iplong_from <= %v AND iplong_to >= %v`

	GetDataCenterV6SqlQuery = `select
	       name,
	       name_code,
	       homepage
	     FROM udger_datacenter_range6
	     JOIN udger_datacenter_list ON udger_datacenter_range6.datacenter_id=udger_datacenter_list.id
           where iplong_from0 <=  %v AND iplong_to0 >= %v AND
		  iplong_from1 <= %v AND iplong_to1 >= %v AND
		  iplong_from2 <= %v AND iplong_to2 >= %v AND
		  iplong_from3 <= %v AND iplong_to3 >= %v AND
		  iplong_from4 <= %v AND iplong_to4 >= %v AND
		  iplong_from5 <= %v AND iplong_to5 >= %v AND
		  iplong_from6 <= %v AND iplong_to6 >= %v AND
		  iplong_from7 <= %v AND iplong_to7 >= %v`

)
