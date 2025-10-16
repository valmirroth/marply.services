package repository

const qrySemestral = `select codcal, codbh, perref, sitafa,
		codccu, dtApuracao, codfil, 
		resumo.CRACHA, resumo.Colaborador, 
		resumo.numcad, resumo.numemp, resumo.tipcol, 
		resumo.ValHoraMes, ValHoraCalculado, resumo.dessit, 
		resumo.codsit, horas, [R$ valor] from (
select codcal, codbh, perref, sitafa,
		codccu, dtApuracao, codfil, 
		resumo.CRACHA, resumo.Colaborador, 
		resumo.numcad, resumo.numemp, resumo.tipcol, 
		resumo.ValHoraMes, ValHoraCalculado, resumo.dessit, 
		resumo.codsit, SUM(qtdhor) as horas, horamin, (ValHoraCalculado * SUM( resumo.qtdhor ) ) as [R$ valor]
		from (
		SELECT ( select case when codbhr = 0 and numcad in ('7259','344756') then 5 else codbhr end from r038hsi where numcad = r034fun.numcad and numemp = r034fun.numemp  and tipcol = r034fun.tipcol
				and datalt = (
							select max(datalt) from r038hsi 
							where 
								numcad = r034fun.numcad 
								and numemp = r034fun.numemp 
								and tipcol = r034fun.tipcol
								and datalt <= r066apu.datapu
							) ) codbh,
			r044cal.codcal, r044cal.perref, r034fun.codccu, convert(varchar,(r066sit.datapu),112) dtApuracao,
			r034fun.codfil, getdate() as DataAtual, 
			cast(r034fun.numcad  as varchar) +' - '+ r034fun.nomfun as Colaborador,
			CAST(r034fun.numcad AS VARCHAR) AS CRACHA,               
			r066apu.numemp,       
			case when datafa >= '20260226' then 1 else r034fun.sitafa end as sitafa,       
			r066apu.numcad,               
			r066apu.datapu,               
			r066apu.tipcol,              
			nomesc AS TURNO ,  
			r034fun.valsal / 
			cast(r034fun.valsal / 
			(hormes/60) as decimal(19,6))as ValHoraMes,
						cast( case when sit.codsit in (331,338) then ( r034fun.valsal / 
			(hormes/60))*1.6
			when sit.codsit in( 330, 337) then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit in (332,339) then ( r034fun.valsal / 
			(hormes/60))*2
			when sit.codsit in (230,109,336) then ( r034fun.valsal / 
			(hormes/60))	  
			when sit.codsit = 313 then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit in ( 302, 340, 333) then (( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5
			when sit.codsit = 301 then (( r034fun.valsal / 
			(hormes/60))*1.5)  
				when sit.codsit = 305 then (( r034fun.valsal / 
			(hormes/60))*2) 
				when sit.codsit = 311 then (( r034fun.valsal / 
			(hormes/60))*1.6) 
				when sit.codsit in(306,342,335 ) then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 2)
				when sit.codsit = 314 then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5)
				when sit.codsit in( 312 , 341, 334) then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.6)

			end as decimal(19,6)) as ValHoraCalculado,
			case when sit.codsit in (109,230,336) then 'Banco de Horas Negativo' else 
			replace(
			replace(
			replace(   
			replace(  
			replace(replace(replace(sit.dessit,'Banco de Horas','B.H.'),'H. Extra','H.E.'),'Horas Extras','H.E.') 
			,'Intrajornada','Intr')
			,'Noturnas','Not')
			,'Intraj','Intr')
			,'Apurad','Apu.')
			end as dessit,    
			case when sit.codsit in (230,109,336) then 230 else sit.codsit end as codsit,
			case when sit.codsit in (230,109,336) then (r066sit.qtdhor/60.000) * (-1) else (r066sit.qtdhor/60.000)  end qtdhor,
			r066sit.qtdhor horamin
		FROM r066sit
			left join r066apu on  r066sit.numemp = r066apu.numemp and r066sit.tipcol = r066apu.tipcol and r066sit.numcad = r066apu.numcad           
			and r066sit.datapu = r066apu.datapu
			left join r034fun (nolock) on r034fun.numcad = r066apu.numcad   
			LEFT JOIN r010sit sit on sit.codsit = r066sit.codsit
			INNER JOIN r006esc (NOLOCK) ON r066apu.codesc = r006esc.codesc
			inner join r010sit on r010sit.codsit = r034fun.sitafa
			inner join r044cal on r044cal.tipcal ='11' and  r066apu.datapu between r044cal.iniapu and r044cal.fimapu and r044cal.numemp = r066apu.numemp 
		where 1=1 
			--and (r010sit.codsit <> '7' or (r034fun.datafa >= '20250826'))
			--and r034fun.codccu in ('cc103','cc789','cc252')
			and convert(varchar,( r066sit.datapu),112) >= '20250926' 
			and convert(varchar,(r066sit.datapu),112) < '20260226'
			and   ( sit.codsit in ( select  codsit from r011eve where codbhr in (5)  )  )
			and codfil in (1,5)
		) as resumo
		where CRACHA <> 6065 
				and numemp = 1
		--	and CRACHA = '4223'
			and codbh = 5
	group by codcal, codbh, perref, codccu, sitafa,
			dtApuracao, codfil, resumo.CRACHA, 
			ValHoraCalculado, Colaborador, resumo.numcad, 
			resumo.numemp,  resumo.tipcol, resumo.ValHoraMes, 
			resumo.dessit, resumo.codsit, horamin

union all 


select codcal, codbh, perref, sitafa,
		codccu, dtApuracao, codfil, 
		resumo.CRACHA, resumo.Colaborador, 
		resumo.numcad, resumo.numemp, resumo.tipcol, 
		resumo.ValHoraMes, ValHoraCalculado, resumo.dessit, 
		resumo.codsit, SUM(qtdhor) as horas, horamin, (ValHoraCalculado * SUM( resumo.qtdhor ) ) as [R$ valor]
		from (
		SELECT ( select codbhr from r038hsi where numcad = r034fun.numcad and numemp = r034fun.numemp  and tipcol = r034fun.tipcol
				and datalt = (
							select max(datalt) from r038hsi 
							where 
								numcad = r034fun.numcad 
								and numemp = r034fun.numemp 
								and tipcol = r034fun.tipcol
								and datalt <= r066apu.datapu
							) ) codbh,
			r044cal.codcal, r044cal.perref, r034fun.codccu, convert(varchar,(r066sit.datapu),112) dtApuracao,
			r034fun.codfil, getdate() as DataAtual, 
			cast(r034fun.numcad  as varchar) +' - '+ r034fun.nomfun as Colaborador,
			CAST(r034fun.numcad AS VARCHAR) AS CRACHA,               
			r066apu.numemp,       
			case when datafa >= '20251126' then 1 else r034fun.sitafa end as sitafa,      
			r066apu.numcad,               
			r066apu.datapu,               
			r066apu.tipcol,              
			nomesc AS TURNO ,  
			r034fun.valsal / 
			cast(r034fun.valsal / 
			(hormes/60) as decimal(19,6))as ValHoraMes,
						cast( case when sit.codsit in (331,338) then ( r034fun.valsal / 
			(hormes/60))*1.6
			when sit.codsit in( 330, 337) then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit in (332,339) then ( r034fun.valsal / 
			(hormes/60))*2
			when sit.codsit in (230,109,336) then ( r034fun.valsal / 
			(hormes/60))	  
			when sit.codsit = 313 then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit in ( 302, 340, 333) then (( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5
			when sit.codsit = 301 then (( r034fun.valsal / 
			(hormes/60))*1.5)  
				when sit.codsit = 305 then (( r034fun.valsal / 
			(hormes/60))*2) 
				when sit.codsit = 311 then (( r034fun.valsal / 
			(hormes/60))*1.6) 
				when sit.codsit in(306,342,335 ) then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 2)
				when sit.codsit = 314 then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5)
				when sit.codsit in( 312 , 341, 334) then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.6)

			end as decimal(19,6)) as ValHoraCalculado,
			case when sit.codsit in (109,230,336) then 'Banco de Horas Negativo' else 
			replace(
			replace(
			replace(   
			replace(  
			replace(replace(replace(sit.dessit,'Banco de Horas','B.H.'),'H. Extra','H.E.'),'Horas Extras','H.E.') 
			,'Intrajornada','Intr')
			,'Noturnas','Not')
			,'Intraj','Intr')
			,'Apurad','Apu.')
			end as dessit,    
			case when sit.codsit in (230,109,336) then 230 else sit.codsit end as codsit,
			case when sit.codsit in (230,109,336) then (r066sit.qtdhor/60.000) * (-1) else (r066sit.qtdhor/60.000)  end qtdhor,
			r066sit.qtdhor horamin
		FROM  (	select numemp, tipcol, numcad, datlan as datapu, codsit,  qtdhor from r011lan where  orilan= 'D' and    r011lan.datcmp >= '20250926') r066sit
			left join r066apu on  r066sit.numemp = r066apu.numemp and r066sit.tipcol = r066apu.tipcol and r066sit.numcad = r066apu.numcad           
			and r066sit.datapu = r066apu.datapu
			left join r034fun (nolock) on r034fun.numcad = r066apu.numcad   
			LEFT JOIN r010sit sit on sit.codsit = r066sit.codsit
			INNER JOIN r006esc (NOLOCK) ON r066apu.codesc = r006esc.codesc
			inner join r010sit on r010sit.codsit = r034fun.sitafa
			inner join r044cal on r044cal.tipcal ='11' and  r066apu.datapu between r044cal.iniapu and r044cal.fimapu and r044cal.numemp = r066apu.numemp 
		where 1=1 
			--and (r010sit.codsit <> '7' or (r034fun.datafa >= '20250826'))
			--and r034fun.codccu in ('cc103','cc789','cc252')
			and convert(varchar,( r066sit.datapu),112) >= '20250926' 
			and convert(varchar,(r066sit.datapu),112) < '20260226'
			and   ( sit.codsit in ( select  codsit from r011eve where codbhr in (5)  )  )
			and codfil in (1,5)
		) as resumo
		where CRACHA <> 6065 
				and numemp = 1
	--	and CRACHA = '4223'
			and codbh =5
	group by codcal, codbh, perref, codccu, sitafa,
			dtApuracao, codfil, resumo.CRACHA, 
			ValHoraCalculado, Colaborador, resumo.numcad, 
			resumo.numemp,  resumo.tipcol, resumo.ValHoraMes, 
			resumo.dessit, resumo.codsit, horamin
) as resumo
order by numemp asc, numcad asc, tipcol asc, dtApuracao asc `

const qryMensal = `select codcal, codbh, perref, sitafa,
		codccu, dtApuracao, codfil, 
		resumo.CRACHA, resumo.Colaborador, 
		resumo.numcad, resumo.numemp, resumo.tipcol, 
		resumo.ValHoraMes, ValHoraCalculado, resumo.dessit, 
		resumo.codsit, horas, [R$ valor] from (
select codcal, codbh, perref, sitafa,
		codccu, dtApuracao, codfil, 
		resumo.CRACHA, resumo.Colaborador, 
		resumo.numcad, resumo.numemp, resumo.tipcol, 
		resumo.ValHoraMes, ValHoraCalculado, resumo.dessit, 
		resumo.codsit, SUM(qtdhor) as horas, horamin, (ValHoraCalculado * SUM( resumo.qtdhor ) ) as [R$ valor]
		from (
		SELECT ( select codbhr from r038hsi where numcad = r034fun.numcad and numemp = r034fun.numemp  and tipcol = r034fun.tipcol
				and datalt = (
							select max(datalt) from r038hsi 
							where 
								numcad = r034fun.numcad 
								and numemp = r034fun.numemp 
								and tipcol = r034fun.tipcol
								and datalt <= r066apu.datapu
							) ) codbh,
			r044cal.codcal, r044cal.perref, r034fun.codccu, convert(varchar,(r066sit.datapu),112) dtApuracao,
			r034fun.codfil, getdate() as DataAtual, 
			cast(r034fun.numcad  as varchar) +' - '+ r034fun.nomfun as Colaborador,
			CAST(r034fun.numcad AS VARCHAR) AS CRACHA,               
			r066apu.numemp,       
			r034fun.sitafa,       
			r066apu.numcad,               
			r066apu.datapu,               
			r066apu.tipcol,              
			nomesc AS TURNO ,  
			r034fun.valsal / 
			cast(r034fun.valsal / 
			(hormes/60) as decimal(19,6))as ValHoraMes,
						cast( case when sit.codsit in (331,338) then ( r034fun.valsal / 
			(hormes/60))*1.6
			when sit.codsit in( 330, 337) then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit in (332,339) then ( r034fun.valsal / 
			(hormes/60))*2
			when sit.codsit in (230,109,336) then ( r034fun.valsal / 
			(hormes/60))	  
			when sit.codsit = 313 then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit in ( 302, 340, 333) then (( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5
			when sit.codsit = 301 then (( r034fun.valsal / 
			(hormes/60))*1.5)  
				when sit.codsit = 305 then (( r034fun.valsal / 
			(hormes/60))*2) 
				when sit.codsit = 311 then (( r034fun.valsal / 
			(hormes/60))*1.6) 
				when sit.codsit in(306,342,335 ) then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 2)
				when sit.codsit = 314 then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5)
				when sit.codsit in( 312 , 341, 334) then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.6)

			end as decimal(19,6)) as ValHoraCalculado,
			case when sit.codsit in (109,230,336) then 'Banco de Horas Negativo' else 
			replace(
			replace(
			replace(   
			replace(  
			replace(replace(replace(sit.dessit,'Banco de Horas','B.H.'),'H. Extra','H.E.'),'Horas Extras','H.E.') 
			,'Intrajornada','Intr')
			,'Noturnas','Not')
			,'Intraj','Intr')
			,'Apurad','Apu.')
			end as dessit,    
			case when sit.codsit in (230,109,336) then 230 else sit.codsit end as codsit,
			case when sit.codsit in (230,109,336) then (r066sit.qtdhor/60.000) * (-1) else (r066sit.qtdhor/60.000)  end qtdhor,
			r066sit.qtdhor horamin
		FROM r066sit
			left join r066apu on  r066sit.numemp = r066apu.numemp and r066sit.tipcol = r066apu.tipcol and r066sit.numcad = r066apu.numcad           
			and r066sit.datapu = r066apu.datapu
			left join r034fun (nolock) on r034fun.numcad = r066apu.numcad   
			LEFT JOIN r010sit sit on sit.codsit = r066sit.codsit
			INNER JOIN r006esc (NOLOCK) ON r066apu.codesc = r006esc.codesc
			inner join r010sit on r010sit.codsit = r034fun.sitafa
			inner join r044cal on r044cal.tipcal ='11' and  r066apu.datapu between r044cal.iniapu and r044cal.fimapu and r044cal.numemp = r066apu.numemp 
		where 1=1 
			--r010sit.codsit <> '7'
			--and r034fun.codccu in ('cc103','cc789','cc252')
			and convert(varchar,( r066sit.datapu),112) >= '20260126' 
			and convert(varchar,(r066sit.datapu),112) < '20260226'
			and   ( sit.codsit in ( select  codsit from r011eve where codbhr in (8)  )  )
			and codfil in (1,5)
		) as resumo
		where CRACHA <> 6065 
				and numemp = 1
		--	and CRACHA = '6708'
			and codbh = 8
	group by codcal, codbh, perref, codccu, sitafa,
			dtApuracao, codfil, resumo.CRACHA, 
			ValHoraCalculado, Colaborador, resumo.numcad, 
			resumo.numemp,  resumo.tipcol, resumo.ValHoraMes, 
			resumo.dessit, resumo.codsit, horamin

union all 


select codcal, codbh, perref, sitafa,
		codccu, dtApuracao, codfil, 
		resumo.CRACHA, resumo.Colaborador, 
		resumo.numcad, resumo.numemp, resumo.tipcol, 
		resumo.ValHoraMes, ValHoraCalculado, resumo.dessit, 
		resumo.codsit, SUM(qtdhor) as horas, horamin, (ValHoraCalculado * SUM( resumo.qtdhor ) ) as [R$ valor]
		from (
		SELECT ( select codbhr from r038hsi where numcad = r034fun.numcad and numemp = r034fun.numemp  and tipcol = r034fun.tipcol
				and datalt = (
							select max(datalt) from r038hsi 
							where 
								numcad = r034fun.numcad 
								and numemp = r034fun.numemp 
								and tipcol = r034fun.tipcol
								and datalt <= r066apu.datapu
							) ) codbh,
			r044cal.codcal, r044cal.perref, r034fun.codccu, convert(varchar,(r066sit.datapu),112) dtApuracao,
			r034fun.codfil, getdate() as DataAtual, 
			cast(r034fun.numcad  as varchar) +' - '+ r034fun.nomfun as Colaborador,
			CAST(r034fun.numcad AS VARCHAR) AS CRACHA,               
			r066apu.numemp,       
			r034fun.sitafa,       
			r066apu.numcad,               
			r066apu.datapu,               
			r066apu.tipcol,              
			nomesc AS TURNO ,  
			r034fun.valsal / 
			cast(r034fun.valsal / 
			(hormes/60) as decimal(19,6))as ValHoraMes,
			cast( case when sit.codsit in (331,338) then ( r034fun.valsal / 
			(hormes/60))*1.6
			when sit.codsit in( 330, 337) then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit in (332,339) then ( r034fun.valsal / 
			(hormes/60))*2
			when sit.codsit in (230,109,336) then ( r034fun.valsal / 
			(hormes/60))	  
			when sit.codsit = 313 then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit in ( 302, 340, 333) then (( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5
			when sit.codsit = 301 then (( r034fun.valsal / 
			(hormes/60))*1.5)  
				when sit.codsit = 305 then (( r034fun.valsal / 
			(hormes/60))*2) 
				when sit.codsit = 311 then (( r034fun.valsal / 
			(hormes/60))*1.6) 
				when sit.codsit in(306,342,335 ) then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 2)
				when sit.codsit = 314 then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5)
				when sit.codsit in( 312 , 341, 334) then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.6)

			end as decimal(19,6)) as ValHoraCalculado,
			case when sit.codsit in (109,230,336) then 'Banco de Horas Negativo' else 
			replace(
			replace(
			replace(   
			replace(  
			replace(replace(replace(sit.dessit,'Banco de Horas','B.H.'),'H. Extra','H.E.'),'Horas Extras','H.E.') 
			,'Intrajornada','Intr')
			,'Noturnas','Not')
			,'Intraj','Intr')
			,'Apurad','Apu.')
			end as dessit,    
			case when sit.codsit in (230,109,336) then 230 else sit.codsit end as codsit,
			case when sit.codsit in (230,109,336) then (r066sit.qtdhor/60.000) * (-1) else (r066sit.qtdhor/60.000)  end qtdhor,
			r066sit.qtdhor horamin
		FROM  (	select numemp, tipcol, numcad, datlan as datapu, codsit,  qtdhor from r011lan where  orilan= 'D' and    r011lan.datcmp >= '20260126') r066sit
			left join r066apu on  r066sit.numemp = r066apu.numemp and r066sit.tipcol = r066apu.tipcol and r066sit.numcad = r066apu.numcad           
			and r066sit.datapu = r066apu.datapu
			left join r034fun (nolock) on r034fun.numcad = r066apu.numcad   
			LEFT JOIN r010sit sit on sit.codsit = r066sit.codsit
			INNER JOIN r006esc (NOLOCK) ON r066apu.codesc = r006esc.codesc
			inner join r010sit on r010sit.codsit = r034fun.sitafa
			inner join r044cal on r044cal.tipcal ='11' and  r066apu.datapu between r044cal.iniapu and r044cal.fimapu and r044cal.numemp = r066apu.numemp 
		where 1=1 
			--r010sit.codsit <> '7'
			--and r034fun.codccu in ('cc103','cc789','cc252')
			and convert(varchar,( r066sit.datapu),112) >= '20260126' 
			and convert(varchar,(r066sit.datapu),112) < '20260226'
			and   ( sit.codsit in ( select  codsit from r011eve where codbhr in (8)  )  )
			and codfil in (1,5)
		) as resumo
		where CRACHA <> 6065 
				and numemp = 1
			--and CRACHA = '6708'
			 and codbh = 8
	group by codcal, codbh, perref, codccu, sitafa,
			dtApuracao, codfil, resumo.CRACHA, 
			ValHoraCalculado, Colaborador, resumo.numcad, 
			resumo.numemp,  resumo.tipcol, resumo.ValHoraMes, 
			resumo.dessit, resumo.codsit, horamin
) as resumo
order by numemp, numcad, tipcol, dtApuracao asc `
