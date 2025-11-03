/**
 * Expression evaluation is no longer performed client-side.
 * When expressions are detected in a spec, server-side validation
 * should be used via the validate() function with an endpoint option.
 * 
 * This module is kept for backwards compatibility but evalExpression
 * will throw an error indicating server-side validation should be used.
 * 
 * @deprecated Use server-side validation via validate() with endpoint option
 */
export function evalExpression(expr: string, obj: Record<string, any>): boolean {
  throw new Error(
    'Expression evaluation is no longer supported client-side. ' +
    'Use validate() with an endpoint option for specs containing expressions (conditions).'
  );
}
