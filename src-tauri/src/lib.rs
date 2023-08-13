use std::ffi::CStr;
use std::ffi::CString;
use std::os::raw::c_char;

#[cfg(target_os = "macos")]
#[link(name = "CoreFoundation", kind = "framework")]
#[link(name = "Security", kind = "framework")]
extern "C" {
    fn Request(uri: *const c_char, data: *const c_char) -> *const c_char;
}

#[cfg(not(target_os = "macos"))]
extern "C" {
    fn Request(uri: *const c_char, data: *const c_char) -> *const c_char;
}

pub fn request(uri: &str, data: &str) -> String {
    let result = unsafe {
        Request(
            CString::new(uri).unwrap().as_ptr(),
            CString::new(data).unwrap().as_ptr(),
        )
    };

    let c_str = unsafe { CStr::from_ptr(result) };
    let string = c_str.to_str().expect("Error translating SQIP from library");

    format!("{}", string)
}
