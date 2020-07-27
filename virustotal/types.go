package virustotal

// VirusTotal holds information about an IP address from https://developers.virustotal.com/v3.0/reference#ip-object
type VirusTotal struct {
	Data struct {
		Attributes struct {
			AsOwner             string `json:"as_owner"`
			Asn                 int    `json:"asn"`
			Continent           string `json:"continent"`
			Country             string `json:"country"`
			LastAnalysisResults struct {
				ADMINUSLabs struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"ADMINUSLabs"`
				AegisLabWebGuard struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"AegisLab WebGuard"`
				AlienVault struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"AlienVault"`
				AntiyAVL struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Antiy-AVL"`
				ArtistsAgainst419 struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Artists Against 419"`
				AutoShun struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"AutoShun"`
				Avira struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Avira"`
				BADWAREINFO struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"BADWARE.INFO"`
				BaiduInternational struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Baidu-International"`
				BitDefender struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"BitDefender"`
				BlockList struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"BlockList"`
				Blueliv struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Blueliv"`
				BotvrijEu struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Botvrij.eu"`
				CINSArmy struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"CINS Army"`
				CLEANMX struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"CLEAN MX"`
				CRDF struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"CRDF"`
				Certego struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Certego"`
				ComodoSiteInspector struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Comodo Site Inspector"`
				ComodoValkyrieVerdict struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Comodo Valkyrie Verdict"`
				CyRadar struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"CyRadar"`
				Cyan struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Cyan"`
				CyberCrime struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"CyberCrime"`
				Cyren struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Cyren"`
				DNS8 struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"DNS8"`
				DrWeb struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Dr.Web"`
				ESET struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"ESET"`
				ESTsecurityThreatInside struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"ESTsecurity-Threat Inside"`
				EmergingThreats struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"EmergingThreats"`
				Emsisoft struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Emsisoft"`
				EonScope struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"EonScope"`
				FeodoTracker struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Feodo Tracker"`
				ForcepointThreatSeeker struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Forcepoint ThreatSeeker"`
				Fortinet struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Fortinet"`
				FraudScore struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"FraudScore"`
				FraudSense struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"FraudSense"`
				GData struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"G-Data"`
				GoogleSafebrowsing struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Google Safebrowsing"`
				GreenSnow struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"GreenSnow"`
				IPsum struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"IPsum"`
				K7AntiVirus struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"K7AntiVirus"`
				Kaspersky struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Kaspersky"`
				Lumu struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Lumu"`
				MalSilo struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"MalSilo"`
				Malc0DeDatabase struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Malc0de Database"`
				Malekal struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Malekal"`
				MalwareDomainBlocklist struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Malware Domain Blocklist"`
				MalwareDomainList struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"MalwareDomainList"`
				MalwarePatrol struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"MalwarePatrol"`
				MalwarebytesHpHosts struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Malwarebytes hpHosts"`
				Malwared struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Malwared"`
				Netcraft struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Netcraft"`
				NotMining struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"NotMining"`
				Nucleon struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Nucleon"`
				OpenPhish struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"OpenPhish"`
				Opera struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Opera"`
				PREBYTES struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"PREBYTES"`
				PhishLabs struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"PhishLabs"`
				PhishingDatabase struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Phishing Database"`
				Phishtank struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Phishtank"`
				QuickHeal struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Quick Heal"`
				Quttera struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Quttera"`
				Rising struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Rising"`
				SCUMWAREOrg struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"SCUMWARE.org"`
				Sangfor struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Sangfor"`
				SecureBrain struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"SecureBrain"`
				Segasec struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Segasec"`
				SnortIPSampleList struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Snort IP sample list"`
				Sophos struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Sophos"`
				Spam404 struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Spam404"`
				Spamhaus struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Spamhaus"`
				StopBadware struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"StopBadware"`
				StopForumSpam struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"StopForumSpam"`
				SucuriSiteCheck struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Sucuri SiteCheck"`
				Tencent struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Tencent"`
				ThreatHive struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"ThreatHive"`
				Threatsourcing struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Threatsourcing"`
				Trustwave struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Trustwave"`
				URLQuery struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"URLQuery"`
				URLhaus struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"URLhaus"`
				VXVault struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"VX Vault"`
				VirusdieExternalSiteScan struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Virusdie External Site Scan"`
				WebSecurityGuard struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Web Security Guard"`
				YandexSafebrowsing struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Yandex Safebrowsing"`
				ZCloudsec struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"ZCloudsec"`
				ZDBZeus struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"ZDB Zeus"`
				ZeroCERT struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"ZeroCERT"`
				Zerofox struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"Zerofox"`
				ZeusTracker struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"ZeusTracker"`
				DesenmascaraMe struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"desenmascara.me"`
				MalwaresComURLChecker struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"malwares.com URL checker"`
				Securolytics struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"securolytics"`
				Zvelo struct {
					Category   string `json:"category"`
					EngineName string `json:"engine_name"`
					Method     string `json:"method"`
					Result     string `json:"result"`
				} `json:"zvelo"`
			} `json:"last_analysis_results"`
			LastAnalysisStats struct {
				Harmless   int `json:"harmless"`
				Malicious  int `json:"malicious"`
				Suspicious int `json:"suspicious"`
				Timeout    int `json:"timeout"`
				Undetected int `json:"undetected"`
			} `json:"last_analysis_stats"`
			LastHTTPSCertificate struct {
				CertSignature struct {
					Signature          string `json:"signature"`
					SignatureAlgorithm string `json:"signature_algorithm"`
				} `json:"cert_signature"`
				Extensions struct {
					One3614111129242       string `json:"1.3.6.1.4.1.11129.2.4.2"`
					CA                     bool   `json:"CA"`
					AuthorityKeyIdentifier struct {
						Keyid string `json:"keyid"`
					} `json:"authority_key_identifier"`
					CaInformationAccess struct {
						CAIssuers string `json:"CA Issuers"`
						OCSP      string `json:"OCSP"`
					} `json:"ca_information_access"`
					CertificatePolicies    []string      `json:"certificate_policies"`
					CrlDistributionPoints  []string      `json:"crl_distribution_points"`
					ExtendedKeyUsage       []string      `json:"extended_key_usage"`
					KeyUsage               []string      `json:"key_usage"`
					SubjectAlternativeName []string      `json:"subject_alternative_name"`
					SubjectKeyIdentifier   string        `json:"subject_key_identifier"`
					Tags                   []interface{} `json:"tags"`
				} `json:"extensions"`
				Issuer struct {
					C  string `json:"C"`
					CN string `json:"CN"`
					O  string `json:"O"`
				} `json:"issuer"`
				PublicKey struct {
					Algorithm string `json:"algorithm"`
					Ec        struct {
						Oid string `json:"oid"`
						Pub string `json:"pub"`
					} `json:"ec"`
				} `json:"public_key"`
				SerialNumber       string `json:"serial_number"`
				SignatureAlgorithm string `json:"signature_algorithm"`
				Size               int    `json:"size"`
				Subject            struct {
					C  string `json:"C"`
					CN string `json:"CN"`
					L  string `json:"L"`
					O  string `json:"O"`
					ST string `json:"ST"`
				} `json:"subject"`
				Tags             []interface{} `json:"tags"`
				Thumbprint       string        `json:"thumbprint"`
				ThumbprintSha256 string        `json:"thumbprint_sha256"`
				Validity         struct {
					NotAfter  string `json:"not_after"`
					NotBefore string `json:"not_before"`
				} `json:"validity"`
				Version string `json:"version"`
			} `json:"last_https_certificate"`
			LastHTTPSCertificateDate int           `json:"last_https_certificate_date"`
			LastModificationDate     int           `json:"last_modification_date"`
			Network                  string        `json:"network"`
			RegionalInternetRegistry string        `json:"regional_internet_registry"`
			Reputation               int           `json:"reputation"`
			Tags                     []interface{} `json:"tags"`
			TotalVotes               struct {
				Harmless  int `json:"harmless"`
				Malicious int `json:"malicious"`
			} `json:"total_votes"`
			WhoisDate int       `json:"whois_date"`
		} `json:"attributes"`
		ID    string `json:"id"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Type string `json:"type"`
	} `json:"data"`
}