package forecast

//2016.6.20
//forecast deviation funtion:
// y=-4*10^(-8)*x^6-2*10^(-5)*x^5+0.0028*x^4-0.1463*x^3+1.9963*x^2+47.613*x+17.273
var forecastAverage []int = []int{
	3600: 0,
	3599: 100,
	3598: 100,
	3597: 200,
	3596: 200,
	3595: 300,
	3594: 400,
	3593: 400,
	3592: 500,
	3591: 500,
	3590: 600,
	3589: 600,
	3588: 700,
	3587: 700,
	3586: 800,
	3585: 800,
	3584: 900,
	3583: 900,
	3582: 900,
	3581: 1000,
	3580: 1000,
	3579: 1000,
	3578: 1100,
	3577: 1100,
	3576: 1100,
	3575: 1100,
	3574: 1100,
	3573: 1100,
	3572: 1200,
	3571: 1200,
	3570: 1200,
	3569: 1200,
	3568: 1200,
	3567: 1200,
	3566: 1200,
	3565: 1200,
	3564: 1200,
	3563: 1200,
	3562: 1200,
	3561: 1200,
	3560: 1300,
	3559: 1300,
	3558: 1300,
	3557: 1300,
	3556: 1300,
	3555: 1300,
	3554: 1300,
	3553: 1300,
	3552: 1300,
	3551: 1400,
	3550: 1400,
	3549: 1400,
	3548: 1500,
	3547: 1500,
	3546: 1500,
	3545: 1500,
	3544: 1500,
	3543: 1500,
	3542: 1500,
	3541: 1500,
	3540: 1500,
}

//20160811
var baseDeviation []int = []int{
	3540: 1300,
	3541: 1300,
	3542: 1300,
	3543: 1300,
	3544: 1300,
	3545: 1300,
	3546: 1300,
	3547: 1300,
	3548: 1300,
	3549: 1300,
	3550: 1300,
	3551: 1300,
	3552: 1300,
	3553: 1300,
	3554: 1300,
	3555: 1300,
	3556: 1300,
	3557: 1200,
	3558: 1200,
	3559: 1200,
	3560: 1200,
	3561: 1200,
	3562: 1100,
	3563: 1100,
	3564: 1100,
	3565: 1000,
	3566: 1000,
	3567: 1000,
	3568: 1000,
	3569: 1000,
	3570: 1000,
	3571: 1000,
	3572: 1000,
	3573: 1000,
	3574: 1000,
	3575: 1000,
	3576: 1000,
	3577: 1000,
	3578: 1000,
	3579: 1000,
	3580: 1000,
	3581: 1000,
	3582: 1000,
	3583: 1000,
	3584: 900,
	3585: 900,
	3586: 800,
	3587: 700,
	3588: 700,
	3589: 700,
	3590: 600,
	3591: 600,
	3592: 600,
	3593: 500,
	3594: 500,
	3595: 400,
	3596: 300,
	3597: 200,
	3598: 100,
	3599: 100,
	3600: 0,
}

var standardLevel []int = []int{
	3540: 0,
	3541: 0,
	3542: 0,
	3543: 0,
	3544: 0,
	3545: 0,
	3546: 0,
	3547: 0,
	3548: 0,
	3549: 0,
	3550: 0,
	3551: 0,
	3552: 0,
	3553: 0,
	3554: 0,
	3555: 0,
	3556: 0,
	3557: 0,
	3558: 100,
	3559: 100,
	3560: 100,
	3561: 100,
	3562: 100,
	3563: 100,
	3564: 100,
	3565: 100,
	3566: 200,
	3567: 100,
	3568: 100,
	3569: 100,
	3570: 100,
	3571: 0,
	3572: 0,
	3573: 0,
	3574: 0,
	3575: 0,
	3576: 0,
	3577: 0,
	3578: 0,
	3579: 0,
	3580: 0,
	3581: 0,
	3582: 0,
	3583: 0,
	3584: 0,
	3585: 100,
	3586: 100,
	3587: 200,
	3588: 300,
	3589: 200,
	3590: 200,
	3591: 200,
	3592: 200,
	3593: 100,
	3594: 200,
	3595: 100,
	3596: 200,
	3597: 300,
	3598: 400,
	3599: 400,
	3600: 400,
}
