use std::time::{Instant};

// Will write the interpreter in Rust, for now, just testing the language

fn main() {
    let start = Instant::now();
    const NOMBRE_DE_NOMBRE: usize = 1_000_000;
    let mut nombres_premiers: Vec<usize> = Vec::with_capacity(NOMBRE_DE_NOMBRE);

    nombres_premiers.push(2);

    for i in (3..=NOMBRE_DE_NOMBRE).step_by(2) {
        let mut premier: bool = true;
        let limit = (i as f64).sqrt() as usize;
        for j in 2..=limit {
            if i%j==0 {
                premier = false;
                break
            }
        }
        if premier {
            nombres_premiers.push(i);
        }
    }
    let duration = start.elapsed();
    println!("Elapsed: {:?}", duration);
}
